package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/nerdgatehub/nerdgate-hub/internal/security"
	_ "modernc.org/sqlite"
)

var domainPattern = regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,63}$`)

type Route struct {
	ID        string    `json:"id"`
	Domain    string    `json:"domain"`
	TargetURL string    `json:"target_url"`
	TLS       bool      `json:"tls"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RouteInput struct {
	Domain    string
	TargetURL string
	TLS       bool
}

type Store struct {
	db *sql.DB
}

func Open(dataDir string) (*Store, error) {
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", filepath.Join(dataDir, "nerdgate.db"))
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)

	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := s.importRoutesJSON(filepath.Join(dataDir, "routes.json")); err != nil {
		_ = db.Close()
		return nil, err
	}

	return s, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) List() []Route {
	routes, err := s.ListRoutes(context.Background())
	if err != nil {
		return nil
	}
	return routes
}

func (s *Store) ListRoutes(ctx context.Context) ([]Route, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, domain, target_url, tls, created_at, updated_at
		FROM routes
		ORDER BY domain ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routes []Route
	for rows.Next() {
		var route Route
		var tls int
		if err := rows.Scan(&route.ID, &route.Domain, &route.TargetURL, &tls, &route.CreatedAt, &route.UpdatedAt); err != nil {
			return nil, err
		}
		route.TLS = tls == 1
		routes = append(routes, route)
	}

	return routes, rows.Err()
}

func (s *Store) Create(input RouteInput) (Route, error) {
	return s.CreateRoute(context.Background(), input)
}

func (s *Store) CreateRoute(ctx context.Context, input RouteInput) (Route, error) {
	if err := validateInput(input); err != nil {
		return Route{}, err
	}

	now := time.Now().UTC()
	route := Route{
		ID:        newID(now),
		Domain:    normalizeDomain(input.Domain),
		TargetURL: strings.TrimSpace(input.TargetURL),
		TLS:       input.TLS,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO routes (id, domain, target_url, tls, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, route.ID, route.Domain, route.TargetURL, boolInt(route.TLS), route.CreatedAt, route.UpdatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return Route{}, fmt.Errorf("domain %q already exists", route.Domain)
		}
		return Route{}, err
	}

	return route, nil
}

func (s *Store) Delete(id string) error {
	return s.DeleteRoute(context.Background(), id)
}

func (s *Store) DeleteRoute(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM routes WHERE id = ?`, id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("route %q not found", id)
	}
	return nil
}

func (s *Store) EnsureAdminUser(ctx context.Context, username, password string) error {
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	return s.ResetPassword(ctx, username, password)
}

func (s *Store) ResetPassword(ctx context.Context, username, password string) error {
	username = strings.TrimSpace(username)
	if username == "" {
		return errors.New("username is required")
	}
	if password == "" {
		return errors.New("password is required")
	}

	hash, err := security.HashPassword(password)
	if err != nil {
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now().UTC()
	if _, err := tx.ExecContext(ctx, `DELETE FROM users WHERE username <> ?`, username); err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO users (username, password_hash, created_at, updated_at)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(username) DO UPDATE SET
			password_hash = excluded.password_hash,
			updated_at = excluded.updated_at
	`, username, hash, now, now)
	if err != nil {
		return err
	}

	secret, err := security.RandomHex(32)
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO settings (key, value, updated_at)
		VALUES (?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET
			value = excluded.value,
			updated_at = excluded.updated_at
	`, "session_secret", secret, time.Now().UTC()); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Store) Authenticate(ctx context.Context, username, password string) (bool, error) {
	var hash string
	err := s.db.QueryRowContext(ctx, `SELECT password_hash FROM users WHERE username = ?`, strings.TrimSpace(username)).Scan(&hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return security.CheckPassword(password, hash), nil
}

func (s *Store) EnsureSessionSecret(ctx context.Context, fallback string) (string, error) {
	value, err := s.Setting(ctx, "session_secret")
	if err == nil && value != "" {
		return value, nil
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}

	if strings.TrimSpace(fallback) != "" {
		if err := s.SetSetting(ctx, "session_secret", fallback); err != nil {
			return "", err
		}
		return fallback, nil
	}

	secret, err := security.RandomHex(32)
	if err != nil {
		return "", err
	}
	if err := s.SetSetting(ctx, "session_secret", secret); err != nil {
		return "", err
	}
	return secret, nil
}

func (s *Store) Setting(ctx context.Context, key string) (string, error) {
	var value string
	err := s.db.QueryRowContext(ctx, `SELECT value FROM settings WHERE key = ?`, key).Scan(&value)
	return value, err
}

func (s *Store) SetSetting(ctx context.Context, key, value string) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO settings (key, value, updated_at)
		VALUES (?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET
			value = excluded.value,
			updated_at = excluded.updated_at
	`, key, value, time.Now().UTC())
	return err
}

func (s *Store) migrate() error {
	_, err := s.db.Exec(`
		PRAGMA journal_mode = WAL;
		PRAGMA foreign_keys = ON;

		CREATE TABLE IF NOT EXISTS routes (
			id TEXT PRIMARY KEY,
			domain TEXT NOT NULL UNIQUE,
			target_url TEXT NOT NULL,
			tls INTEGER NOT NULL DEFAULT 1,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);

		CREATE TABLE IF NOT EXISTS users (
			username TEXT PRIMARY KEY,
			password_hash TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);

		CREATE TABLE IF NOT EXISTS settings (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);
	`)
	return err
}

func (s *Store) importRoutesJSON(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	if len(data) == 0 {
		return nil
	}

	var count int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM routes`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	var routes []Route
	if err := json.Unmarshal(data, &routes); err != nil {
		return err
	}

	for _, route := range routes {
		if route.ID == "" {
			route.ID = newID(time.Now().UTC())
		}
		if route.CreatedAt.IsZero() {
			route.CreatedAt = time.Now().UTC()
		}
		if route.UpdatedAt.IsZero() {
			route.UpdatedAt = route.CreatedAt
		}
		_, err := s.db.Exec(`
			INSERT OR IGNORE INTO routes (id, domain, target_url, tls, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`, route.ID, normalizeDomain(route.Domain), route.TargetURL, boolInt(route.TLS), route.CreatedAt, route.UpdatedAt)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateInput(input RouteInput) error {
	domain := normalizeDomain(input.Domain)
	if !domainPattern.MatchString(domain) {
		return errors.New("domain must be a valid hostname, for example app.example.com")
	}

	target := strings.TrimSpace(input.TargetURL)
	parsed, err := url.Parse(target)
	if err != nil {
		return fmt.Errorf("target URL is invalid: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return errors.New("target URL must start with http:// or https://")
	}
	if parsed.Host == "" {
		return errors.New("target URL must include a host")
	}
	if parsed.Port() != "" {
		if _, err := net.LookupPort("tcp", parsed.Port()); err != nil {
			return fmt.Errorf("target URL port is invalid: %w", err)
		}
	}

	return nil
}

func normalizeDomain(domain string) string {
	return strings.ToLower(strings.TrimSpace(domain))
}

func newID(now time.Time) string {
	return fmt.Sprintf("route-%d", now.UnixNano())
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
