package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/nerdgatehub/nerdgate-hub/internal/config"
	"github.com/nerdgatehub/nerdgate-hub/internal/dockerclient"
	"github.com/nerdgatehub/nerdgate-hub/internal/store"
	"github.com/nerdgatehub/nerdgate-hub/internal/traefik"
	"github.com/nerdgatehub/nerdgate-hub/internal/web"
)

func main() {
	cfg := config.FromEnv()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	routeStore, err := store.Open(filepath.Join(cfg.DataDir, "routes.json"))
	if err != nil {
		logger.Error("open store", "error", err)
		os.Exit(1)
	}

	renderer := traefik.NewRenderer(cfg.TraefikDynamicPath)
	if err := renderer.Render(routeStore.List()); err != nil {
		logger.Error("render traefik config", "error", err)
		os.Exit(1)
	}

	app := web.NewServer(web.ServerConfig{
		Username:      cfg.Username,
		Password:      cfg.Password,
		SessionSecret: cfg.SessionSecret,
		Store:         routeStore,
		Renderer:      renderer,
		Docker:        dockerclient.New(cfg.DockerSocketPath, cfg.DockerProxyNetwork),
		Logger:        logger,
	})

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           app.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("nerdgate hub listening", "addr", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("http server failed", "error", err)
			os.Exit(1)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown failed", "error", err)
		os.Exit(1)
	}
}
