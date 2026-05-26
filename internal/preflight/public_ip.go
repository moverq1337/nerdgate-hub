package preflight

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/netip"
	"strings"
	"time"
)

var publicIPEndpoints = []string{
	"https://api.ipify.org",
	"https://ifconfig.me/ip",
	"https://icanhazip.com",
}

func DetectPublicIP(ctx context.Context) (netip.Addr, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	var lastErr error

	for _, endpoint := range publicIPEndpoints {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
		if err != nil {
			lastErr = err
			continue
		}

		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		body, readErr := io.ReadAll(io.LimitReader(resp.Body, 128))
		closeErr := resp.Body.Close()
		if readErr != nil {
			lastErr = readErr
			continue
		}
		if closeErr != nil {
			lastErr = closeErr
			continue
		}
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			lastErr = fmt.Errorf("%s returned %s", endpoint, resp.Status)
			continue
		}

		addr, err := netip.ParseAddr(strings.TrimSpace(string(body)))
		if err != nil {
			lastErr = err
			continue
		}
		return addr.Unmap(), nil
	}

	if lastErr != nil {
		return netip.Addr{}, lastErr
	}
	return netip.Addr{}, fmt.Errorf("no public IP endpoints configured")
}
