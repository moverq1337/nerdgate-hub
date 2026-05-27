package web

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/nerdgatehub/nerdgate-hub/internal/dockerclient"
	"github.com/nerdgatehub/nerdgate-hub/internal/store"
)

type DiagnosticsData struct {
	ConfigPath   string
	ConfigStatus string
	BackupPath   string
	RestoreHint  string
	DockerStatus string
	HealthStatus string
	TraefikLogs  string
	AuditEvents  []store.AuditEvent
	Hints        []string
}

func (s *Server) diagnostics(ctx context.Context, routes []RouteView, containers []dockerclient.Container, dockerErr error) DiagnosticsData {
	data := DiagnosticsData{
		ConfigPath:  s.renderer.Path(),
		BackupPath:  "/backup",
		RestoreHint: "Use scripts/restore.sh on the host to restore a backup safely.",
		Hints:       make([]string, 0),
	}

	if info, err := os.Stat(data.ConfigPath); err != nil {
		data.ConfigStatus = "Missing: " + err.Error()
		data.Hints = append(data.Hints, "Traefik dynamic config is missing. Try saving a route or restarting NerdGate Hub.")
	} else {
		data.ConfigStatus = fmt.Sprintf("OK, updated %s", info.ModTime().Format("2006-01-02 15:04:05 MST"))
	}

	if dockerErr != nil {
		data.DockerStatus = "Unavailable: " + dockerErr.Error()
		data.Hints = append(data.Hints, "Docker API is unavailable. Container picker, attach action, and Traefik log diagnostics need /var/run/docker.sock.")
	} else {
		data.DockerStatus = fmt.Sprintf("OK, %d running containers visible", len(containers))
	}

	up, down := healthCounts(routes)
	data.HealthStatus = fmt.Sprintf("%d up, %d down", up, down)
	if down > 0 {
		data.Hints = append(data.Hints, "One or more targets are down. Open the health chip for details, then check the target URL, Docker network, firewall, and backend logs.")
	}

	if dockerErr == nil {
		if container, ok := findTraefikContainer(containers); ok {
			logs, err := s.docker.ContainerLogs(ctx, container.ID, 120)
			if err != nil {
				data.TraefikLogs = "Could not read Traefik logs: " + err.Error()
			} else {
				data.TraefikLogs = logs
				data.Hints = append(data.Hints, traefikHints(logs)...)
			}
		} else {
			data.TraefikLogs = "Traefik container was not found in Docker API results."
			data.Hints = append(data.Hints, "Traefik container is not visible. Check docker compose ps on the host.")
		}
	}

	if len(data.Hints) == 0 {
		data.Hints = append(data.Hints, "No known certificate or routing problem detected from the current checks.")
	}
	if events, err := s.store.ListAuditEvents(ctx, 12); err == nil {
		data.AuditEvents = events
	} else {
		s.logger.Warn("list audit events failed", "error", err)
	}

	return data
}

func healthCounts(routes []RouteView) (int, int) {
	up := 0
	down := 0
	for _, route := range routes {
		switch route.Health.State {
		case "up":
			up++
		case "down":
			down++
		}
	}
	return up, down
}

func findTraefikContainer(containers []dockerclient.Container) (dockerclient.Container, bool) {
	for _, container := range containers {
		if container.Labels["com.docker.compose.service"] == "traefik" {
			return container, true
		}
	}
	for _, container := range containers {
		if strings.Contains(strings.ToLower(container.Name), "traefik") {
			return container, true
		}
	}
	return dockerclient.Container{}, false
}

func traefikHints(logs string) []string {
	lower := strings.ToLower(logs)
	hints := make([]string, 0)

	if strings.Contains(lower, "unable to obtain acme certificate") ||
		strings.Contains(lower, "error getting certificate") ||
		strings.Contains(lower, "acme: error") {
		hints = append(hints, "Let's Encrypt could not issue a certificate. Run check-domain for the domain and verify public TCP 80/443 reach Traefik.")
	}
	if strings.Contains(lower, "connection refused") ||
		strings.Contains(lower, "timeout") ||
		strings.Contains(lower, "timeout during connect") {
		hints = append(hints, "HTTP-01 validation likely cannot reach Traefik. Open TCP 80/443 and disable CDN proxying while issuing.")
	}
	if strings.Contains(lower, "ratelimited") ||
		strings.Contains(lower, "too many certificates") ||
		strings.Contains(lower, "too many failed authorizations") {
		hints = append(hints, "Let's Encrypt rate limit was hit. Fix DNS/ports first, then wait for the limit window.")
	}
	if strings.Contains(lower, "permission denied") && strings.Contains(lower, "acme") {
		hints = append(hints, "ACME storage permissions look wrong. On the host, run chmod 600 data/acme/acme.json.")
	}
	if strings.Contains(lower, "address already in use") {
		hints = append(hints, "Port 80 or 443 is already occupied. Stop another proxy or move it away from public ports.")
	}

	return hints
}
