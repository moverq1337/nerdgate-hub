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

type setupPageData struct {
	Error string
}

func (s *Server) setup(w http.ResponseWriter, r *http.Request) {
	hasUsers, err := s.store.HasUsers(r.Context())
	if err != nil {
		s.logger.Error("check setup state", "error", err)
		http.Error(w, "setup state failed", http.StatusInternalServerError)
		return
	}
	if hasUsers {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := setupPageData{
		Error: r.URL.Query().Get("error"),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.tmpl.ExecuteTemplate(w, "setup.html", data); err != nil {
		s.logger.Error("render setup", "error", err)
		http.Error(w, "render failed", http.StatusInternalServerError)
	}
}

func (s *Server) setupPost(w http.ResponseWriter, r *http.Request) {
	hasUsers, err := s.store.HasUsers(r.Context())
	if err != nil {
		s.logger.Error("check setup state", "error", err)
		http.Redirect(w, r, "/setup?error="+url.QueryEscape("Setup failed"), http.StatusSeeOther)
		return
	}
	if hasUsers {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/setup?error="+url.QueryEscape("Invalid form"), http.StatusSeeOther)
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	password := r.FormValue("password")
	if password != r.FormValue("confirm_password") {
		http.Redirect(w, r, "/setup?error="+url.QueryEscape("Passwords do not match"), http.StatusSeeOther)
		return
	}

	if err := s.store.CompleteSetup(r.Context(), r.FormValue("setup_token"), username, password); err != nil {
		s.logger.Warn("setup failed", "error", err)
		http.Redirect(w, r, "/setup?error="+url.QueryEscape(err.Error()), http.StatusSeeOther)
		return
	}

	if secret, err := s.store.Setting(r.Context(), "session_secret"); err == nil && secret != "" {
		s.setSessionSecret(secret)
	} else if err != nil {
		s.logger.Warn("setup completed but session secret refresh failed", "error", err)
	}

	http.SetCookie(w, s.newSessionCookie(username))
	http.Redirect(w, r, "/?status=setup-complete", http.StatusSeeOther)
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	hasUsers, err := s.store.HasUsers(r.Context())
	if err != nil {
		s.logger.Error("check setup state", "error", err)
		http.Error(w, "setup state failed", http.StatusInternalServerError)
		return
	}
	if !hasUsers {
		http.Redirect(w, r, "/setup", http.StatusSeeOther)
		return
	}

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
	hasUsers, err := s.store.HasUsers(r.Context())
	if err != nil {
		s.logger.Error("check setup state", "error", err)
		http.Redirect(w, r, "/login?error="+url.QueryEscape("Login failed"), http.StatusSeeOther)
		return
	}
	if !hasUsers {
		http.Redirect(w, r, "/setup", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/login?error="+url.QueryEscape("Invalid form"), http.StatusSeeOther)
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	ok, err := s.store.Authenticate(r.Context(), username, r.FormValue("password"))
	if err != nil {
		s.logger.Error("authenticate user", "error", err)
		http.Redirect(w, r, "/login?error="+url.QueryEscape("Login failed"), http.StatusSeeOther)
		return
	}
	if !ok {
		http.Redirect(w, r, "/login?error="+url.QueryEscape("Invalid login or password"), http.StatusSeeOther)
		return
	}

	http.SetCookie(w, s.newSessionCookie(username))
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
		hasUsers, err := s.store.HasUsers(r.Context())
		if err != nil {
			s.logger.Error("check setup state", "error", err)
			http.Error(w, "setup state failed", http.StatusInternalServerError)
			return
		}
		if !hasUsers {
			http.Redirect(w, r, "/setup", http.StatusSeeOther)
			return
		}

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
	if !ok {
		return false
	}
	authenticated, err := s.store.Authenticate(r.Context(), user, pass)
	if err != nil {
		s.logger.Warn("basic auth failed", "error", err)
		return false
	}
	return authenticated
}

func (s *Server) newSessionCookie(username string) *http.Cookie {
	expires := time.Now().Add(30 * 24 * time.Hour)
	payload := fmt.Sprintf("%s:%d", username, expires.Unix())
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
	if len(parts) != 2 || parts[0] == "" {
		return false
	}

	expires, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return false
	}

	return time.Now().Unix() < expires
}

func (s *Server) sign(payload string) string {
	secret := s.currentSessionSecret()
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(payload))
	return encodeBytes(mac.Sum(nil))
}

func (s *Server) setSessionSecret(value string) {
	s.sessionMu.Lock()
	defer s.sessionMu.Unlock()
	s.sessionSecret = sessionSecret(value)
}

func (s *Server) currentSessionSecret() []byte {
	s.sessionMu.RLock()
	defer s.sessionMu.RUnlock()
	return append([]byte(nil), s.sessionSecret...)
}

func sessionSecret(value string) []byte {
	sum := sha256.Sum256([]byte(value))
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
