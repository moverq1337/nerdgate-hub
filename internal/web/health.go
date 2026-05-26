package web

import (
	"context"
	"fmt"
	"net/http"
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
		Detail: err.Error(),
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
			Detail: fmt.Sprintf("HTTP %d", status),
			Class:  "down",
		}
	}
}
