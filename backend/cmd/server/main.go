package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ustkost/korotkiy-url/internal/config"
	"github.com/ustkost/korotkiy-url/internal/db"
	"github.com/ustkost/korotkiy-url/internal/handler"
	"github.com/ustkost/korotkiy-url/internal/repository"
	"github.com/ustkost/korotkiy-url/internal/service"
)

func main() {
	if err := run(); err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run() error {

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.LogLevel,
	})))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := db.Connect(ctx, cfg)
	if err != nil {
		return err
	}
	defer pool.Close()

	linkRepo := repository.NewLinkRepository(pool)
	clickRepo := repository.NewClickRepository(pool)

	linkService := service.NewLinkService(linkRepo)
	clickService := service.NewClickService(clickRepo)

	linkHandler := handler.NewLinkHandler(linkService)
	clickHandler := handler.NewClickHandler(clickService)
	redirectHandler := handler.NewRedirectHandler(linkService, clickService)

	mux := handler.NewRouter(linkHandler, clickHandler, redirectHandler)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handler.Logging(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		slog.Info("listening", "port", cfg.Port)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		slog.Info("shutting down...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
