package web

import (
	"html/template"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/nerdgatehub/nerdgate-hub/internal/dockerclient"
	"github.com/nerdgatehub/nerdgate-hub/internal/store"
	"github.com/nerdgatehub/nerdgate-hub/internal/traefik"
)

type ServerConfig struct {
	SessionSecret string
	Store         *store.Store
	Renderer      *traefik.Renderer
	Docker        *dockerclient.Client
	Logger        *slog.Logger
}

type Server struct {
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

func NewServer(cfg ServerConfig) *Server {
	return &Server{
		sessionSecret: sessionSecret(cfg.SessionSecret),
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

func (s *Server) redirectError(w http.ResponseWriter, r *http.Request, message string) {
	message = strings.TrimSpace(message)
	if message == "" {
		message = "request failed"
	}
	http.Redirect(w, r, "/?error="+url.QueryEscape(message), http.StatusSeeOther)
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
