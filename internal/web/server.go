package web

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/nerdgatehub/nerdgate-hub/internal/dockerclient"
	"github.com/nerdgatehub/nerdgate-hub/internal/store"
	"github.com/nerdgatehub/nerdgate-hub/internal/traefik"
)

type ServerConfig struct {
	Username      string
	Password      string
	SessionSecret string
	Store         *store.Store
	Renderer      *traefik.Renderer
	Docker        *dockerclient.Client
	Logger        *slog.Logger
}

type Server struct {
	username      string
	password      string
	sessionSecret []byte
	store         *store.Store
	renderer      *traefik.Renderer
	docker        *dockerclient.Client
	logger        *slog.Logger
	tmpl          *template.Template
}

type pageData struct {
	Routes                 []store.Route
	DockerTargets          []dockerclient.TargetOption
	AvailableDockerTargets int
	DockerError            string
	Status                 string
	Error                  string
}

type loginPageData struct {
	Error string
	Next  string
}

func NewServer(cfg ServerConfig) *Server {
	return &Server{
		username:      cfg.Username,
		password:      cfg.Password,
		sessionSecret: sessionSecret(cfg.Username, cfg.Password, cfg.SessionSecret),
		store:         cfg.Store,
		renderer:      cfg.Renderer,
		docker:        cfg.Docker,
		logger:        cfg.Logger,
		tmpl:          template.Must(template.ParseFS(templates, "templates/*.html")),
	}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.health)
	mux.Handle("GET /static/", http.FileServerFS(staticFiles))
	mux.HandleFunc("GET /login", s.login)
	mux.HandleFunc("POST /login", s.loginPost)
	mux.HandleFunc("POST /logout", s.withAuth(s.logout))
	mux.HandleFunc("GET /", s.withAuth(s.index))
	mux.HandleFunc("POST /routes", s.withAuth(s.createRoute))
	mux.HandleFunc("POST /routes/{id}/delete", s.withAuth(s.deleteRoute))
	return mux
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	dockerTargets, dockerErr := s.dockerTargets(r)

	data := pageData{
		Routes:                 s.store.List(),
		DockerTargets:          dockerTargets,
		AvailableDockerTargets: availableTargetCount(dockerTargets),
		Status:                 r.URL.Query().Get("status"),
		Error:                  r.URL.Query().Get("error"),
	}
	if dockerErr != nil {
		s.logger.Warn("docker discovery failed", "error", dockerErr)
		data.DockerError = "Docker socket is unavailable. Manual target URL still works."
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		s.logger.Error("render index", "error", err)
		http.Error(w, "render failed", http.StatusInternalServerError)
	}
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	if s.isAuthenticated(r) {
		http.Redirect(w, r, safeNext(r.URL.Query().Get("next")), http.StatusSeeOther)
		return
	}

	data := loginPageData{
		Error: r.URL.Query().Get("error"),
		Next:  safeNext(r.URL.Query().Get("next")),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.tmpl.ExecuteTemplate(w, "login.html", data); err != nil {
		s.logger.Error("render login", "error", err)
		http.Error(w, "render failed", http.StatusInternalServerError)
	}
}

func (s *Server) loginPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/login?error="+url.QueryEscape("Invalid form"), http.StatusSeeOther)
		return
	}

	if !constantTimeEqual(r.FormValue("username"), s.username) || !constantTimeEqual(r.FormValue("password"), s.password) {
		http.Redirect(w, r, "/login?error="+url.QueryEscape("Invalid login or password"), http.StatusSeeOther)
		return
	}

	http.SetCookie(w, s.newSessionCookie())
	http.Redirect(w, r, safeNext(r.FormValue("next")), http.StatusSeeOther)
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "nerdgate_session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (s *Server) createRoute(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.redirectError(w, r, "invalid form")
		return
	}

	input := store.RouteInput{
		Domain:    r.FormValue("domain"),
		TargetURL: routeTargetURL(r),
		TLS:       r.FormValue("tls") == "on",
	}

	if _, err := s.store.Create(input); err != nil {
		s.redirectError(w, r, err.Error())
		return
	}

	if err := s.renderer.Render(s.store.List()); err != nil {
		s.logger.Error("render traefik config", "error", err)
		s.redirectError(w, r, "route saved, but Traefik config render failed")
		return
	}

	http.Redirect(w, r, "/?status=created", http.StatusSeeOther)
}

func (s *Server) deleteRoute(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.store.Delete(id); err != nil {
		s.redirectError(w, r, err.Error())
		return
	}

	if err := s.renderer.Render(s.store.List()); err != nil {
		s.logger.Error("render traefik config", "error", err)
		s.redirectError(w, r, "route deleted, but Traefik config render failed")
		return
	}

	http.Redirect(w, r, "/?status=deleted", http.StatusSeeOther)
}

func (s *Server) dockerTargets(r *http.Request) ([]dockerclient.TargetOption, error) {
	if s.docker == nil {
		return nil, nil
	}

	containers, err := s.docker.ListContainers(r.Context())
	if err != nil {
		return nil, err
	}

	return dockerclient.TargetOptions(containers), nil
}

func routeTargetURL(r *http.Request) string {
	if r.FormValue("target_mode") == "docker" {
		return r.FormValue("docker_target")
	}
	return r.FormValue("target_url")
}

func (s *Server) withAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.isAuthenticated(r) {
			nextURL := url.QueryEscape(r.URL.RequestURI())
			http.Redirect(w, r, "/login?next="+nextURL, http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

func (s *Server) redirectError(w http.ResponseWriter, r *http.Request, message string) {
	message = strings.TrimSpace(message)
	if message == "" {
		message = "request failed"
	}
	http.Redirect(w, r, "/?error="+url.QueryEscape(message), http.StatusSeeOther)
}

func constantTimeEqual(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

func (s *Server) isAuthenticated(r *http.Request) bool {
	if s.validSession(r) {
		return true
	}

	user, pass, ok := r.BasicAuth()
	return ok && constantTimeEqual(user, s.username) && constantTimeEqual(pass, s.password)
}

func (s *Server) newSessionCookie() *http.Cookie {
	expires := time.Now().Add(30 * 24 * time.Hour)
	payload := fmt.Sprintf("%s:%d", s.username, expires.Unix())
	return &http.Cookie{
		Name:     "nerdgate_session",
		Value:    encode(payload) + "." + s.sign(payload),
		Path:     "/",
		Expires:  expires,
		MaxAge:   int((30 * 24 * time.Hour).Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func (s *Server) validSession(r *http.Request) bool {
	cookie, err := r.Cookie("nerdgate_session")
	if err != nil {
		return false
	}

	encodedPayload, signature, ok := strings.Cut(cookie.Value, ".")
	if !ok {
		return false
	}

	payload, err := base64.RawURLEncoding.DecodeString(encodedPayload)
	if err != nil {
		return false
	}

	if !constantTimeEqual(signature, s.sign(string(payload))) {
		return false
	}

	parts := strings.Split(string(payload), ":")
	if len(parts) != 2 || parts[0] != s.username {
		return false
	}

	expires, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return false
	}

	return time.Now().Unix() < expires
}

func (s *Server) sign(payload string) string {
	mac := hmac.New(sha256.New, s.sessionSecret)
	_, _ = mac.Write([]byte(payload))
	return encodeBytes(mac.Sum(nil))
}

func sessionSecret(username, password, configured string) []byte {
	seed := configured
	if seed == "" {
		seed = username + ":" + password
	}
	sum := sha256.Sum256([]byte(seed))
	return sum[:]
}

func encode(value string) string {
	return encodeBytes([]byte(value))
}

func encodeBytes(value []byte) string {
	return base64.RawURLEncoding.EncodeToString(value)
}

func safeNext(value string) string {
	if value == "" || !strings.HasPrefix(value, "/") || strings.HasPrefix(value, "//") {
		return "/"
	}
	return value
}

func availableTargetCount(options []dockerclient.TargetOption) int {
	count := 0
	for _, option := range options {
		if option.Available {
			count++
		}
	}
	return count
}
