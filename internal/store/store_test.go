package store

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestStoreRoutesAndAuth(t *testing.T) {
	ctx := context.Background()
	s, err := Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	if err := s.EnsureAdminUser(ctx, "admin", "secret"); err != nil {
		t.Fatal(err)
	}

	ok, err := s.Authenticate(ctx, "admin", "secret")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected admin login to work")
	}

	if _, err := s.CreateRoute(ctx, RouteInput{
		Domain:    "App.Example.com",
		TargetURL: "http://host.docker.internal:3000",
		TLS:       true,
	}); err != nil {
		t.Fatal(err)
	}

	routes, err := s.ListRoutes(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(routes) != 1 {
		t.Fatalf("expected 1 route, got %d", len(routes))
	}
	if routes[0].Domain != "app.example.com" {
		t.Fatalf("unexpected normalized domain: %s", routes[0].Domain)
	}
}

func TestResetPasswordReplacesPreviousAdmin(t *testing.T) {
	ctx := context.Background()
	s, err := Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	if err := s.ResetPassword(ctx, "admin", "old-secret"); err != nil {
		t.Fatal(err)
	}
	if err := s.ResetPassword(ctx, "owner", "new-secret"); err != nil {
		t.Fatal(err)
	}

	ok, err := s.Authenticate(ctx, "admin", "old-secret")
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected previous admin to be removed")
	}

	ok, err = s.Authenticate(ctx, "owner", "new-secret")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected new admin to work")
	}
}

func TestStoreImportsRoutesJSON(t *testing.T) {
	dir := t.TempDir()
	data := `[
  {
    "id": "route-existing",
    "domain": "legacy.example.com",
    "target_url": "http://host.docker.internal:8080",
    "tls": true,
    "created_at": "2026-05-27T00:00:00Z",
    "updated_at": "2026-05-27T00:00:00Z"
  }
]`
	if err := os.WriteFile(filepath.Join(dir, "routes.json"), []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}

	s, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	routes, err := s.ListRoutes(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(routes) != 1 {
		t.Fatalf("expected imported route, got %d", len(routes))
	}
	if routes[0].Domain != "legacy.example.com" {
		t.Fatalf("unexpected imported domain: %s", routes[0].Domain)
	}
}
