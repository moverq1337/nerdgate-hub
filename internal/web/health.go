package web

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/nerdgatehub/nerdgate-hub/internal/store"
)

type RouteView struct {
	Route  store.Route
	Health HealthStatus
}

type HealthStatus struct {
	State  string
	Label  string
	Detail string
	Class  string
}

func routeViews(ctx context.Context, routes []store.Route) []RouteView {
	views := make([]RouteView, len(routes))
	for i, route := range routes {
		views[i] = RouteView{
			Route: route,
			Health: HealthStatus{
				State: "checking",
				Label: "Checking",
				Class: "checking",
			},
		}
	}

	var wg sync.WaitGroup
	limit := make(chan struct{}, 8)
	for i := range views {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			limit <- struct{}{}
			defer func() { <-limit }()
			views[index].Health = checkTarget(ctx, views[index].Route.TargetURL)
		}(i)
	}
	wg.Wait()

	return views
}

func checkTarget(parent context.Context, targetURL string) HealthStatus {
	ctx, cancel := context.WithTimeout(parent, 1800*time.Millisecond)
	defer cancel()

	client := &http.Client{Timeout: 1800 * time.Millisecond}
	status, err := requestTarget(ctx, client, http.MethodHead, targetURL)
	if err == nil && status != http.StatusMethodNotAllowed {
		return healthFromStatus(status)
	}
	if err == nil && status == http.StatusMethodNotAllowed {
		status, err = requestTarget(ctx, client, http.MethodGet, targetURL)
		if err == nil {
			return healthFromStatus(status)
		}
	}

	return HealthStatus{
		State:  "down",
		Label:  "Down",
		Detail: friendlyTargetError(err),
		Class:  "down",
	}
}

func requestTarget(ctx context.Context, client *http.Client, method, targetURL string) (int, error) {
	req, err := http.NewRequestWithContext(ctx, method, targetURL, nil)
	if err != nil {
		return 0, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

func healthFromStatus(status int) HealthStatus {
	switch {
	case status >= 200 && status < 500:
		return HealthStatus{
			State:  "up",
			Label:  "Up",
			Detail: fmt.Sprintf("HTTP %d", status),
			Class:  "up",
		}
	default:
		return HealthStatus{
			State:  "down",
			Label:  "Down",
			Detail: fmt.Sprintf("Target returned HTTP %d. Check the backend logs and target URL.", status),
			Class:  "down",
		}
	}
}

func friendlyTargetError(err error) string {
	if err == nil {
		return "Target did not respond."
	}
	value := strings.ToLower(err.Error())
	switch {
	case strings.Contains(value, "connection refused"):
		return "Connection refused. The service is not listening on that host/port, or the container is not attached to nerdgate-proxy."
	case strings.Contains(value, "no such host"):
		return "Host was not found. Check the container name, DNS name, or target URL."
	case strings.Contains(value, "i/o timeout"), strings.Contains(value, "context deadline exceeded"):
		return "Target timed out. Check firewall rules, Docker network access, and whether the backend is running."
	case strings.Contains(value, "certificate"):
		return "Target TLS failed. Use http:// for internal backends unless the backend has a valid certificate."
	default:
		return err.Error()
	}
}
