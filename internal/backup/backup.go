package backup

import (
	"archive/zip"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type Metadata struct {
	CreatedAt string `json:"created_at"`
	Version   string `json:"version"`
}

type Info struct {
	HasDatabase bool
	HasRoutes   bool
	HasACME     bool
	HasMetadata bool
}

const pendingRestoreName = "restore-pending.zip"

func Create(ctx context.Context, dataDir, acmePath, outputPath string) error {
	dataDir = filepath.Clean(dataDir)
	if outputPath == "" {
		return fmt.Errorf("backup output path is required")
	}
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return err
	}

	tmpDir, err := os.MkdirTemp("", "nerdgate-backup-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(dataDir, "nerdgate.db")
	snapshotPath := filepath.Join(tmpDir, "nerdgate.db")
	if err := snapshotSQLite(ctx, dbPath, snapshotPath); err != nil {
		return err
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	archive := zip.NewWriter(file)
	if err := addFile(archive, snapshotPath, "nerdgate.db"); err != nil {
		archive.Close()
		return err
	}

	if exists(filepath.Join(dataDir, "routes.json")) {
		if err := addFile(archive, filepath.Join(dataDir, "routes.json"), "routes.json"); err != nil {
			archive.Close()
			return err
		}
	}
	if strings.TrimSpace(acmePath) != "" && exists(acmePath) {
		if err := addFile(archive, acmePath, "acme.json"); err != nil {
			archive.Close()
			return err
		}
	}

	meta, err := json.MarshalIndent(Metadata{
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		Version:   "1",
	}, "", "  ")
	if err != nil {
		archive.Close()
		return err
	}
	if err := addBytes(archive, "metadata.json", meta); err != nil {
		archive.Close()
		return err
	}

	return archive.Close()
}

func Inspect(inputPath string) (Info, error) {
	reader, err := zip.OpenReader(inputPath)
	if err != nil {
		return Info{}, err
	}
	defer reader.Close()

	var info Info
	for _, file := range reader.File {
		switch file.Name {
		case "nerdgate.db":
			info.HasDatabase = true
		case "routes.json":
			info.HasRoutes = true
		case "acme.json":
			info.HasACME = true
		case "metadata.json":
			info.HasMetadata = true
		}
	}
	if !info.HasDatabase {
		return info, fmt.Errorf("backup does not contain nerdgate.db")
	}
	return info, nil
}

func Restore(dataDir, acmePath, inputPath string) error {
	if _, err := Inspect(inputPath); err != nil {
		return err
	}

	reader, err := zip.OpenReader(inputPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return err
	}
	if strings.TrimSpace(acmePath) != "" {
		if err := os.MkdirAll(filepath.Dir(acmePath), 0o755); err != nil {
			return err
		}
	}

	restoredDB := false
	for _, file := range reader.File {
		switch file.Name {
		case "nerdgate.db":
			if err := extractFile(file, filepath.Join(dataDir, "nerdgate.db"), 0o600); err != nil {
				return err
			}
			restoredDB = true
			_ = os.Remove(filepath.Join(dataDir, "nerdgate.db-wal"))
			_ = os.Remove(filepath.Join(dataDir, "nerdgate.db-shm"))
		case "routes.json":
			if err := extractFile(file, filepath.Join(dataDir, "routes.json"), 0o600); err != nil {
				return err
			}
		case "acme.json":
			if strings.TrimSpace(acmePath) != "" {
				if err := extractFile(file, acmePath, 0o600); err != nil {
					return err
				}
			}
		}
	}
	if !restoredDB {
		return fmt.Errorf("backup does not contain nerdgate.db")
	}

	return nil
}

func PendingRestorePath(dataDir string) string {
	return filepath.Join(filepath.Clean(dataDir), pendingRestoreName)
}

func StageRestore(dataDir, inputPath string) (Info, error) {
	info, err := Inspect(inputPath)
	if err != nil {
		return info, err
	}
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return info, err
	}

	outputPath := PendingRestorePath(dataDir)
	tmpPath := outputPath + ".tmp"
	if err := copyFile(inputPath, tmpPath, 0o600); err != nil {
		return info, err
	}
	if err := os.Rename(tmpPath, outputPath); err != nil {
		_ = os.Remove(tmpPath)
		return info, err
	}
	return info, nil
}

func ApplyPendingRestore(dataDir, acmePath string) (bool, error) {
	path := PendingRestorePath(dataDir)
	if !exists(path) {
		return false, nil
	}
	if err := Restore(dataDir, acmePath, path); err != nil {
		failedPath := path + ".failed-" + time.Now().UTC().Format("20060102-150405")
		_ = os.Rename(path, failedPath)
		return true, err
	}
	if err := os.Remove(path); err != nil {
		return true, err
	}
	return true, nil
}

func snapshotSQLite(ctx context.Context, dbPath, outputPath string) error {
	if !exists(dbPath) {
		return fmt.Errorf("database not found: %s", dbPath)
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.ExecContext(ctx, "VACUUM INTO "+sqliteQuote(outputPath))
	return err
}

func addFile(archive *zip.Writer, path, name string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer, err := archive.Create(name)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, file)
	return err
}

func addBytes(archive *zip.Writer, name string, value []byte) error {
	writer, err := archive.Create(name)
	if err != nil {
		return err
	}
	_, err = writer.Write(value)
	return err
}

func extractFile(file *zip.File, outputPath string, mode os.FileMode) error {
	reader, err := file.Open()
	if err != nil {
		return err
	}
	defer reader.Close()

	output, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = io.Copy(output, reader)
	return err
}

func copyFile(inputPath, outputPath string, mode os.FileMode) error {
	input, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = io.Copy(output, input)
	return err
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func sqliteQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "''") + "'"
}
