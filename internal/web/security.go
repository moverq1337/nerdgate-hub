package web

import (
	"crypto/rand"
	"encoding/hex"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

const csrfCookieName = "nerdgate_csrf"

type RateLimiter struct {
	mu       sync.Mutex
	limit    int
	window   time.Duration
	attempts map[string]rateBucket
}

type rateBucket struct {
	count int
	reset time.Time
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		limit:    limit,
		window:   window,
		attempts: make(map[string]rateBucket),
	}
}

func (r *RateLimiter) Allow(key string) bool {
	now := time.Now()
	r.mu.Lock()
	defer r.mu.Unlock()

	bucket := r.attempts[key]
	if now.After(bucket.reset) {
		bucket = rateBucket{reset: now.Add(r.window)}
	}
	bucket.count++
	r.attempts[key] = bucket
	return bucket.count <= r.limit
}

func (r *RateLimiter) Reset(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.attempts, key)
}

func (s *Server) csrfToken(w http.ResponseWriter, r *http.Request) string {
	if cookie, err := r.Cookie(csrfCookieName); err == nil && cookie.Value != "" {
		return cookie.Value
	}

	token := randomToken()
	http.SetCookie(w, &http.Cookie{
		Name:     csrfCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   int((12 * time.Hour).Seconds()),
		HttpOnly: true,
		Secure:   requestIsSecure(r),
		SameSite: http.SameSiteLaxMode,
	})
	return token
}

func (s *Server) validCSRF(r *http.Request) bool {
	cookie, err := r.Cookie(csrfCookieName)
	if err != nil || cookie.Value == "" {
		return false
	}
	return constantTimeEqual(cookie.Value, r.FormValue("_csrf"))
}

func randomToken() string {
	var value [32]byte
	if _, err := rand.Read(value[:]); err != nil {
		return hex.EncodeToString([]byte(time.Now().Format(time.RFC3339Nano)))
	}
	return hex.EncodeToString(value[:])
}

func requestIsSecure(r *http.Request) bool {
	return r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https")
}

func clientIP(r *http.Request) string {
	if forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); forwarded != "" {
		first, _, _ := strings.Cut(forwarded, ",")
		if ip := strings.TrimSpace(first); ip != "" {
			return ip
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && host != "" {
		return host
	}
	return r.RemoteAddr
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "same-origin")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; img-src 'self' data:; style-src 'self'; script-src 'self'; base-uri 'self'; frame-ancestors 'none'")
		next.ServeHTTP(w, r)
	})
}
