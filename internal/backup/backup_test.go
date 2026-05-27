package backup_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/nerdgatehub/nerdgate-hub/internal/backup"
	"github.com/nerdgatehub/nerdgate-hub/internal/store"
)

func TestCreateAndRestore(t *testing.T) {
	ctx := context.Background()
	sourceDir := t.TempDir()
	sourceAcme := filepath.Join(t.TempDir(), "acme.json")

	db, err := store.Open(sourceDir)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.EnsureAdminUser(ctx, "admin", "secret"); err != nil {
		t.Fatal(err)
	}
	if _, err := db.CreateRoute(ctx, store.RouteInput{
		Domain:    "app.example.com",
		TargetURL: "http://host.docker.internal:3000",
		TLS:       true,
	}); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(sourceAcme, []byte(`{"Account":{}}`), 0o600); err != nil {
		t.Fatal(err)
	}

	archive := filepath.Join(t.TempDir(), "backup.zip")
	if err := backup.Create(ctx, sourceDir, sourceAcme, archive); err != nil {
		t.Fatal(err)
	}
	info, err := backup.Inspect(archive)
	if err != nil {
		t.Fatal(err)
	}
	if !info.HasDatabase || !info.HasACME || !info.HasMetadata {
		t.Fatalf("unexpected backup info: %#v", info)
	}

	restoreDir := t.TempDir()
	restoreAcme := filepath.Join(t.TempDir(), "acme.json")
	if err := backup.Restore(restoreDir, restoreAcme, archive); err != nil {
		t.Fatal(err)
	}

	restored, err := store.Open(restoreDir)
	if err != nil {
		t.Fatal(err)
	}
	defer restored.Close()

	ok, err := restored.Authenticate(ctx, "admin", "secret")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected restored admin credentials")
	}

	routes, err := restored.ListRoutes(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(routes) != 1 || routes[0].Domain != "app.example.com" {
		t.Fatalf("unexpected restored routes: %#v", routes)
	}

	acme, err := os.ReadFile(restoreAcme)
	if err != nil {
		t.Fatal(err)
	}
	if string(acme) != `{"Account":{}}` {
		t.Fatalf("unexpected restored acme content: %s", acme)
	}
}

func TestStageAndApplyPendingRestore(t *testing.T) {
	ctx := context.Background()
	sourceDir := t.TempDir()

	db, err := store.Open(sourceDir)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.EnsureAdminUser(ctx, "admin", "secret"); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	archive := filepath.Join(t.TempDir(), "backup.zip")
	if err := backup.Create(ctx, sourceDir, "", archive); err != nil {
		t.Fatal(err)
	}

	restoreDir := t.TempDir()
	if _, err := backup.StageRestore(restoreDir, archive); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(backup.PendingRestorePath(restoreDir)); err != nil {
		t.Fatal(err)
	}

	applied, err := backup.ApplyPendingRestore(restoreDir, "")
	if err != nil {
		t.Fatal(err)
	}
	if !applied {
		t.Fatal("expected pending restore to apply")
	}
	if _, err := os.Stat(backup.PendingRestorePath(restoreDir)); !os.IsNotExist(err) {
		t.Fatalf("expected pending restore file to be removed, got %v", err)
	}

	restored, err := store.Open(restoreDir)
	if err != nil {
		t.Fatal(err)
	}
	defer restored.Close()

	ok, err := restored.Authenticate(ctx, "admin", "secret")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected restored admin credentials")
	}
}
