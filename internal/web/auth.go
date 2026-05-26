package web

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type loginPageData struct {
	Error string
	Next  string
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

func constantTimeEqual(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
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
