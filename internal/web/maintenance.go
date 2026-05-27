package web

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nerdgatehub/nerdgate-hub/internal/backup"
)

const maxRestoreUploadBytes = 25 << 20

func (s *Server) changePassword(w http.ResponseWriter, r *http.Request) {
	if !s.validCSRF(r) {
		s.redirectError(w, r, "Security token expired. Reload the page and try again.")
		return
	}
	if err := r.ParseForm(); err != nil {
		s.redirectError(w, r, "Invalid password form.")
		return
	}

	actor := currentActor(r)
	if actor == "" {
		s.redirectError(w, r, "Could not determine the current user. Log out and sign in again.")
		return
	}

	currentPassword := r.FormValue("current_password")
	newPassword := r.FormValue("new_password")
	if newPassword != r.FormValue("confirm_password") {
		s.redirectError(w, r, "New passwords do not match.")
		return
	}
	if len(newPassword) < 8 {
		s.redirectError(w, r, "New password must be at least 8 characters.")
		return
	}

	ok, err := s.store.Authenticate(r.Context(), actor, currentPassword)
	if err != nil {
		s.logger.Error("authenticate before password change", "error", err)
		s.redirectError(w, r, "Could not verify the current password.")
		return
	}
	if !ok {
		s.audit(r.Context(), "password.change_failed", actor, clientIP(r))
		s.redirectError(w, r, "Current password is incorrect.")
		return
	}

	if err := s.store.ResetPassword(r.Context(), actor, newPassword); err != nil {
		s.logger.Error("change password", "error", err)
		s.redirectError(w, r, "Could not update the password.")
		return
	}
	if secret, err := s.store.Setting(r.Context(), "session_secret"); err == nil && secret != "" {
		s.setSessionSecret(secret)
	} else if err != nil {
		s.logger.Warn("password updated but session secret refresh failed", "error", err)
	}

	s.audit(r.Context(), "password.change", actor, clientIP(r))
	http.SetCookie(w, s.newSessionCookie(r, actor))
	http.Redirect(w, r, "/?status=password-updated", http.StatusSeeOther)
}

func (s *Server) stageRestore(w http.ResponseWriter, r *http.Request) {
	if !s.validCSRF(r) {
		s.redirectError(w, r, "Security token expired. Reload the page and try again.")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxRestoreUploadBytes)
	if err := r.ParseMultipartForm(maxRestoreUploadBytes); err != nil {
		s.redirectError(w, r, "Restore upload is invalid or larger than 25 MB.")
		return
	}
	if r.MultipartForm != nil {
		defer r.MultipartForm.RemoveAll()
	}

	if r.FormValue("confirm_restore") != "on" {
		s.redirectError(w, r, "Confirm that you understand restore will restart NerdGate Hub.")
		return
	}

	file, header, err := r.FormFile("backup")
	if err != nil {
		s.redirectError(w, r, "Choose a NerdGate backup .zip file.")
		return
	}
	defer file.Close()

	if !strings.HasSuffix(strings.ToLower(header.Filename), ".zip") {
		s.redirectError(w, r, "Restore file must be a .zip backup.")
		return
	}

	tmpDir, err := os.MkdirTemp("", "nerdgate-restore-upload-*")
	if err != nil {
		s.logger.Error("create restore temp dir", "error", err)
		s.redirectError(w, r, "Could not prepare restore upload.")
		return
	}
	defer os.RemoveAll(tmpDir)

	tmpPath := filepath.Join(tmpDir, "restore.zip")
	output, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		s.logger.Error("create restore upload file", "error", err)
		s.redirectError(w, r, "Could not save restore upload.")
		return
	}
	if _, err := io.Copy(output, file); err != nil {
		output.Close()
		s.logger.Error("write restore upload", "error", err)
		s.redirectError(w, r, "Could not save restore upload.")
		return
	}
	if err := output.Close(); err != nil {
		s.logger.Error("close restore upload", "error", err)
		s.redirectError(w, r, "Could not save restore upload.")
		return
	}

	info, err := backup.StageRestore(s.dataDir, tmpPath)
	if err != nil {
		s.logger.Warn("stage restore failed", "error", err)
		s.redirectError(w, r, "Backup restore failed validation: "+err.Error())
		return
	}

	detail := fmt.Sprintf("%s db=%t routes=%t acme=%t metadata=%t", header.Filename, info.HasDatabase, info.HasRoutes, info.HasACME, info.HasMetadata)
	s.audit(r.Context(), "backup.restore_staged", currentActor(r), detail)

	go func() {
		time.Sleep(900 * time.Millisecond)
		s.logger.Info("exiting to apply staged restore")
		os.Exit(0)
	}()

	http.Redirect(w, r, "/?status=restore-staged", http.StatusSeeOther)
}
