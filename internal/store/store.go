package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
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
	mu     sync.Mutex
	path   string
	routes []Route
}

func Open(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}

	s := &Store{path: path}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return s, nil
		}
		return nil, err
	}
	if len(data) == 0 {
		return s, nil
	}
	if err := json.Unmarshal(data, &s.routes); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Store) List() []Route {
	s.mu.Lock()
	defer s.mu.Unlock()

	routes := append([]Route(nil), s.routes...)
	sort.Slice(routes, func(i, j int) bool {
		return routes[i].Domain < routes[j].Domain
	})
	return routes
}

func (s *Store) Create(input RouteInput) (Route, error) {
	if err := validateInput(input); err != nil {
		return Route{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	domain := normalizeDomain(input.Domain)
	for _, route := range s.routes {
		if route.Domain == domain {
			return Route{}, fmt.Errorf("domain %q already exists", domain)
		}
	}

	now := time.Now().UTC()
	route := Route{
		ID:        newID(now),
		Domain:    domain,
		TargetURL: strings.TrimSpace(input.TargetURL),
		TLS:       input.TLS,
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.routes = append(s.routes, route)
	if err := s.saveLocked(); err != nil {
		return Route{}, err
	}

	return route, nil
}

func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, route := range s.routes {
		if route.ID == id {
			s.routes = append(s.routes[:i], s.routes[i+1:]...)
			return s.saveLocked()
		}
	}

	return fmt.Errorf("route %q not found", id)
}

func (s *Store) saveLocked() error {
	data, err := json.MarshalIndent(s.routes, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
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
