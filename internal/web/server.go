package web

import (
	"context"
	"html/template"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/nerdgatehub/nerdgate-hub/internal/backup"
	"github.com/nerdgatehub/nerdgate-hub/internal/dockerclient"
	"github.com/nerdgatehub/nerdgate-hub/internal/store"
	"github.com/nerdgatehub/nerdgate-hub/internal/traefik"
)

type ServerConfig struct {
	SessionSecret string
	DataDir       string
	AcmePath      string
	Store         *store.Store
	Renderer      *traefik.Renderer
	Docker        *dockerclient.Client
	Logger        *slog.Logger
}

type Server struct {
	sessionMu     sync.RWMutex
	sessionSecret []byte
	dataDir       string
	acmePath      string
	store         *store.Store
	renderer      *traefik.Renderer
	docker        *dockerclient.Client
	logger        *slog.Logger
	tmpl          *template.Template
	rateLimiter   *RateLimiter
}

type pageData struct {
	Routes                 []RouteView
	DockerTargets          []dockerclient.TargetOption
	AttachableTargets      []dockerclient.TargetOption
	AvailableDockerTargets int
	Diagnostics            DiagnosticsData
	CSRFToken              string
	DockerError            string
	StatusMessage          string
	Error                  string
}

var domainSplitPattern = regexp.MustCompile(`[,\s]+`)

func NewServer(cfg ServerConfig) *Server {
	return &Server{
		sessionSecret: sessionSecret(cfg.SessionSecret),
		dataDir:       cfg.DataDir,
		acmePath:      cfg.AcmePath,
		store:         cfg.Store,
		renderer:      cfg.Renderer,
		docker:        cfg.Docker,
		logger:        cfg.Logger,
		tmpl:          template.Must(template.ParseFS(templates, "templates/*.html")),
		rateLimiter:   NewRateLimiter(8, 10*time.Minute),
	}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.health)
	mux.Handle("GET /static/", http.FileServerFS(staticFiles))
	mux.HandleFunc("GET /setup", s.setup)
	mux.HandleFunc("POST /setup", s.setupPost)
	mux.HandleFunc("GET /login", s.login)
	mux.HandleFunc("POST /login", s.loginPost)
	mux.HandleFunc("POST /logout", s.withAuth(s.logout))
	mux.HandleFunc("GET /", s.withAuth(s.index))
	mux.HandleFunc("GET /backup", s.withAuth(s.downloadBackup))
	mux.HandleFunc("POST /restore", s.withAuth(s.stageRestore))
	mux.HandleFunc("POST /account/password", s.withAuth(s.changePassword))
	mux.HandleFunc("POST /routes", s.withAuth(s.createRoute))
	mux.HandleFunc("POST /routes/{id}", s.withAuth(s.updateRoute))
	mux.HandleFunc("POST /routes/{id}/delete", s.withAuth(s.deleteRoute))
	mux.HandleFunc("POST /containers/{id}/attach", s.withAuth(s.attachContainer))
	return securityHeaders(mux)
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) downloadBackup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpDir, err := os.MkdirTemp("", "nerdgate-backup-download-*")
	if err != nil {
		http.Error(w, "backup failed", http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tmpDir)

	name := "nerdgate-backup-" + time.Now().UTC().Format("20060102-150405") + ".zip"
	path := filepath.Join(tmpDir, name)
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	if err := backup.Create(ctx, s.dataDir, s.acmePath, path); err != nil {
		s.logger.Error("create backup", "error", err)
		http.Error(w, "backup failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	s.audit(r.Context(), "backup.download", currentActor(r), name)

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", `attachment; filename="`+name+`"`)
	http.ServeFile(w, r, path)
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	routes := s.store.List()
	routeViews := routeViews(r.Context(), routes)
	containers, dockerErr := s.dockerContainers(r)
	dockerTargets := dockerclient.TargetOptions(containers)

	data := pageData{
		Routes:                 routeViews,
		DockerTargets:          dockerTargets,
		AttachableTargets:      attachableTargets(dockerTargets),
		AvailableDockerTargets: availableTargetCount(dockerTargets),
		Diagnostics:            s.diagnostics(r.Context(), routeViews, containers, dockerErr),
		CSRFToken:              s.csrfToken(w, r),
		StatusMessage:          statusMessage(r.URL.Query().Get("status")),
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
	if !s.validCSRF(r) {
		s.redirectError(w, r, "security token expired; reload the page and try again")
		return
	}
	if err := r.ParseForm(); err != nil {
		s.redirectError(w, r, "invalid form")
		return
	}

	domains := parseDomains(r.FormValue("domains"))
	if len(domains) == 0 {
		s.redirectError(w, r, "at least one domain is required")
		return
	}

	inputs := make([]store.RouteInput, 0, len(domains))
	for _, domain := range domains {
		inputs = append(inputs, store.RouteInput{
			Domain:    domain,
			TargetURL: routeTargetURL(r),
			TLS:       r.FormValue("tls") == "on",
		})
	}

	if _, err := s.store.CreateRoutes(r.Context(), inputs); err != nil {
		s.redirectError(w, r, err.Error())
		return
	}
	s.audit(r.Context(), "route.create", currentActor(r), strings.Join(domains, ", "))

	if err := s.renderer.Render(s.store.List()); err != nil {
		s.logger.Error("render traefik config", "error", err)
		s.redirectError(w, r, "route saved, but Traefik config render failed")
		return
	}

	http.Redirect(w, r, "/?status=created", http.StatusSeeOther)
}

func (s *Server) updateRoute(w http.ResponseWriter, r *http.Request) {
	if !s.validCSRF(r) {
		s.redirectError(w, r, "security token expired; reload the page and try again")
		return
	}
	if err := r.ParseForm(); err != nil {
		s.redirectError(w, r, "invalid form")
		return
	}

	input := store.RouteInput{
		Domain:    r.FormValue("domain"),
		TargetURL: r.FormValue("target_url"),
		TLS:       r.FormValue("tls") == "on",
	}

	if _, err := s.store.Update(r.PathValue("id"), input); err != nil {
		s.redirectError(w, r, err.Error())
		return
	}
	s.audit(r.Context(), "route.update", currentActor(r), input.Domain)

	if err := s.renderer.Render(s.store.List()); err != nil {
		s.logger.Error("render traefik config", "error", err)
		s.redirectError(w, r, "route updated, but Traefik config render failed")
		return
	}

	http.Redirect(w, r, "/?status=updated", http.StatusSeeOther)
}

func (s *Server) deleteRoute(w http.ResponseWriter, r *http.Request) {
	if !s.validCSRF(r) {
		s.redirectError(w, r, "security token expired; reload the page and try again")
		return
	}
	id := r.PathValue("id")
	if err := s.store.Delete(id); err != nil {
		s.redirectError(w, r, err.Error())
		return
	}
	s.audit(r.Context(), "route.delete", currentActor(r), id)

	if err := s.renderer.Render(s.store.List()); err != nil {
		s.logger.Error("render traefik config", "error", err)
		s.redirectError(w, r, "route deleted, but Traefik config render failed")
		return
	}

	http.Redirect(w, r, "/?status=deleted", http.StatusSeeOther)
}

func (s *Server) attachContainer(w http.ResponseWriter, r *http.Request) {
	if !s.validCSRF(r) {
		s.redirectError(w, r, "security token expired; reload the page and try again")
		return
	}
	if s.docker == nil {
		s.redirectError(w, r, "Docker API is unavailable")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := s.docker.ConnectContainerToProxyNetwork(ctx, r.PathValue("id")); err != nil {
		s.logger.Warn("attach container failed", "error", err)
		s.redirectError(w, r, err.Error())
		return
	}
	s.audit(r.Context(), "container.attach", currentActor(r), r.PathValue("id"))

	http.Redirect(w, r, "/?status=attached", http.StatusSeeOther)
}

func (s *Server) audit(ctx context.Context, action, actor, detail string) {
	if err := s.store.AddAuditEvent(ctx, action, actor, detail); err != nil {
		s.logger.Warn("audit event failed", "action", action, "error", err)
	}
}

func (s *Server) dockerContainers(r *http.Request) ([]dockerclient.Container, error) {
	if s.docker == nil {
		return nil, nil
	}

	return s.docker.ListContainers(r.Context())
}

func routeTargetURL(r *http.Request) string {
	if r.FormValue("target_mode") == "docker" {
		return r.FormValue("docker_target")
	}
	return r.FormValue("target_url")
}

func (s *Server) redirectError(w http.ResponseWriter, r *http.Request, message string) {
	message = friendlyUserError(message)
	if message == "" {
		message = "request failed"
	}
	http.Redirect(w, r, "/?error="+url.QueryEscape(message), http.StatusSeeOther)
}

func friendlyUserError(message string) string {
	message = strings.TrimSpace(message)
	lower := strings.ToLower(message)
	switch {
	case strings.Contains(lower, "security token"):
		return "Security token expired. Reload the page and try again."
	case strings.Contains(lower, "domain must be a valid hostname"):
		return "Domain must be a hostname like app.example.com. Do not include http://, paths, or ports."
	case strings.Contains(lower, "target url must start"):
		return "Target URL must start with http:// or https://."
	case strings.Contains(lower, "target url must include"):
		return "Target URL must include a host, for example http://203.0.113.10:8080."
	case strings.Contains(lower, "already exists"):
		return message
	case strings.Contains(lower, "docker api"):
		return "Docker API is unavailable. Check that /var/run/docker.sock is mounted into NerdGate Hub."
	default:
		return message
	}
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

func attachableTargets(options []dockerclient.TargetOption) []dockerclient.TargetOption {
	targets := make([]dockerclient.TargetOption, 0)
	seen := make(map[string]bool)
	for _, option := range options {
		if option.Available || option.Kind != "unavailable" || option.ContainerID == "" {
			continue
		}
		key := option.ContainerID
		if seen[key] {
			continue
		}
		seen[key] = true
		targets = append(targets, option)
	}
	return targets
}

func parseDomains(value string) []string {
	parts := domainSplitPattern.Split(value, -1)
	domains := make([]string, 0, len(parts))
	seen := make(map[string]bool)
	for _, part := range parts {
		domain := strings.ToLower(strings.TrimSpace(part))
		if domain == "" || seen[domain] {
			continue
		}
		seen[domain] = true
		domains = append(domains, domain)
	}
	return domains
}

func statusMessage(status string) string {
	switch status {
	case "created":
		return "Route created."
	case "deleted":
		return "Route deleted."
	case "updated":
		return "Route updated."
	case "attached":
		return "Container attached to nerdgate-proxy. Internal targets are available now."
	case "setup-complete":
		return "Setup complete. Welcome to NerdGate Hub."
	case "password-updated":
		return "Password updated. Existing sessions were rotated."
	case "restore-staged":
		return "Backup accepted. NerdGate Hub is restarting to apply the restore."
	default:
		return ""
	}
}
