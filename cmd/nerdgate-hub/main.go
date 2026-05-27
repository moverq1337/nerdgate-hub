package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nerdgatehub/nerdgate-hub/internal/config"
	"github.com/nerdgatehub/nerdgate-hub/internal/dockerclient"
	"github.com/nerdgatehub/nerdgate-hub/internal/store"
	"github.com/nerdgatehub/nerdgate-hub/internal/traefik"
	"github.com/nerdgatehub/nerdgate-hub/internal/web"
)

func main() {
	if len(os.Args) > 1 {
		runCommand(os.Args[1], os.Args[2:])
		return
	}

	cfg := config.FromEnv()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	routeStore, err := store.Open(cfg.DataDir)
	if err != nil {
		logger.Error("open store", "error", err)
		os.Exit(1)
	}
	defer routeStore.Close()

	setupCtx := context.Background()
	hasUsers, err := routeStore.HasUsers(setupCtx)
	if err != nil {
		logger.Error("check admin users", "error", err)
		os.Exit(1)
	}
	if !hasUsers {
		switch {
		case cfg.SetupToken != "":
			if err := routeStore.EnsureSetupToken(setupCtx, cfg.SetupToken); err != nil {
				logger.Error("ensure setup token", "error", err)
				os.Exit(1)
			}
			logger.Info("first-run setup is waiting for browser setup")
		case cfg.Password != "":
			if err := routeStore.EnsureAdminUser(setupCtx, cfg.Username, cfg.Password); err != nil {
				logger.Error("ensure admin user", "error", err)
				os.Exit(1)
			}
		default:
			logger.Error("first-run setup token is required", "hint", "set NERDGATE_SETUP_TOKEN or NERDGATE_PASSWORD")
			os.Exit(1)
		}
	}

	sessionSecret, err := routeStore.EnsureSessionSecret(setupCtx, cfg.SessionSecret)
	if err != nil {
		logger.Error("ensure session secret", "error", err)
		os.Exit(1)
	}

	renderer := traefik.NewRenderer(cfg.TraefikDynamicPath)
	if err := renderer.Render(routeStore.List()); err != nil {
		logger.Error("render traefik config", "error", err)
		os.Exit(1)
	}

	app := web.NewServer(web.ServerConfig{
		SessionSecret: sessionSecret,
		DataDir:       cfg.DataDir,
		AcmePath:      cfg.AcmePath,
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
