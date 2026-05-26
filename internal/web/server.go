package web

import (
	"crypto/subtle"
	"html/template"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/nerdgatehub/nerdgate-hub/internal/store"
	"github.com/nerdgatehub/nerdgate-hub/internal/traefik"
)

type ServerConfig struct {
	Username string
	Password string
	Store    *store.Store
	Renderer *traefik.Renderer
	Logger   *slog.Logger
}

type Server struct {
	username string
	password string
	store    *store.Store
	renderer *traefik.Renderer
	logger   *slog.Logger
	tmpl     *template.Template
}

type pageData struct {
	Routes []store.Route
	Status string
	Error  string
}

func NewServer(cfg ServerConfig) *Server {
	return &Server{
		username: cfg.Username,
		password: cfg.Password,
		store:    cfg.Store,
		renderer: cfg.Renderer,
		logger:   cfg.Logger,
		tmpl:     template.Must(template.ParseFS(templates, "templates/*.html")),
	}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.health)
	mux.Handle("GET /static/app.css", http.FileServerFS(staticFiles))
	mux.HandleFunc("GET /", s.withAuth(s.index))
	mux.HandleFunc("POST /routes", s.withAuth(s.createRoute))
	mux.HandleFunc("POST /routes/{id}/delete", s.withAuth(s.deleteRoute))
	return mux
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	data := pageData{
		Routes: s.store.List(),
		Status: r.URL.Query().Get("status"),
		Error:  r.URL.Query().Get("error"),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		s.logger.Error("render index", "error", err)
		http.Error(w, "render failed", http.StatusInternalServerError)
	}
}

func (s *Server) createRoute(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.redirectError(w, r, "invalid form")
		return
	}

	input := store.RouteInput{
		Domain:    r.FormValue("domain"),
		TargetURL: r.FormValue("target_url"),
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

func (s *Server) withAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || !constantTimeEqual(user, s.username) || !constantTimeEqual(pass, s.password) {
			w.Header().Set("WWW-Authenticate", `Basic realm="NerdGate Hub", charset="UTF-8"`)
			http.Error(w, "authentication required", http.StatusUnauthorized)
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
