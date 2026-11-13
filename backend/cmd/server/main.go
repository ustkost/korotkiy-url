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
	"sync/atomic"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ustkost/korotkiy-url/internal/config"
	"github.com/ustkost/korotkiy-url/internal/db"
	"github.com/ustkost/korotkiy-url/internal/handler"
	"github.com/ustkost/korotkiy-url/internal/metrics"
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
	setupLogger(cfg.LogLevel)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done();
		stop()
	}()

	pool, err := db.Connect(ctx, cfg)
	if err != nil {
		return err
	}
	defer pool.Close()

	return New(ctx, cfg, pool).Run(ctx)
}

func setupLogger(level slog.Level) {
	slogHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if len(groups) == 0 {
				switch a.Key {
				case slog.TimeKey:
					a.Key = "timestamp"
				case slog.MessageKey:
					a.Key = "message"
				}
			}
			return a
		},
	})
	slog.SetDefault(slog.New(slogHandler).With("service.name", "korotkiy-url-backend"))
}

type App struct {
	ready      atomic.Bool
	apiSrv     *http.Server
	metricsSrv *http.Server
}

func New(ctx context.Context, cfg *config.Config, pool *pgxpool.Pool) *App {
	a := &App{}

	linkRepo := repository.NewLinkRepository(pool)
	clickRepo := repository.NewClickRepository(pool)

	linkService := service.NewLinkService(linkRepo)
	clickService := service.NewClickService(clickRepo)

	health := handler.NewHealthHandler(a.ready.Load, pool.Ping)
	linkHandler := handler.NewLinkHandler(linkService)
	clickHandler := handler.NewClickHandler(clickService)
	redirectHandler := handler.NewRedirectHandler(linkService, clickService)

	metrics.RegisterLinksGauge(ctx, linkRepo.Count)
	metrics.RedirectsTotal.WithLabelValues("found").Add(0)
	metrics.RedirectsTotal.WithLabelValues("not_found").Add(0)
	metrics.RedirectsTotal.WithLabelValues("error").Add(0)

	mux := handler.NewRouter(health, linkHandler, clickHandler, redirectHandler)

	handler.InitChaos()
	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.Handler())
	metricsMux.HandleFunc("/admin/chaos", handler.ChaosHandler)

	a.apiSrv = &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handler.Metrics(handler.Logging(mux)),
		ReadHeaderTimeout: 5 * time.Second,
	}
	a.metricsSrv = &http.Server{
		Addr:              ":" + cfg.MetricsPort,
		Handler:           metricsMux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	return a
}

func (a *App) Run(ctx context.Context) error {
	a.ready.Store(true)

	errCh := make(chan error, 2)
	go func() {
		slog.Info("listening api", "addr", a.apiSrv.Addr)
		errCh <- a.apiSrv.ListenAndServe()
	}()
	go func() {
		slog.Info("listening metrics", "addr", a.metricsSrv.Addr)
		errCh <- a.metricsSrv.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		slog.Info("shutting down...")
	}

	a.ready.Store(false)
	time.Sleep(5 * time.Second)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return errors.Join(
		a.apiSrv.Shutdown(shutdownCtx),
		a.metricsSrv.Shutdown(shutdownCtx),
	)
}
