package metrics

import (
	"context"
	"log/slog"
	"math"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	LinksCreated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "korotkiy_url_links_created_total",
		Help: "Total number of created short links",
	})

	RedirectsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "korotkiy_url_redirects_total",
		Help: "Total number of redirects by result",
	}, []string{"result"})

	RedirectDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "korotkiy_url_redirect_duration_seconds",
		Help: "Redirect duration",
	})

	HttpRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"method", "route", "code"})

	HttpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request duration in seconds",
		Buckets: []float64{0.0005, 0.001, 0.002, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1},
	}, []string{"method", "route"})

	HttpInFlight = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "http_requests_in_flight",
		Help: "Number of HTTP requests currently being served",
	})
)

func RegisterLinksGauge(ctx context.Context, count func(ctx context.Context) (int64, error)) {
	promauto.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "korotkiy_url_links_total",
		Help: "Current number of links",
	}, func() float64 {
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		n, err := count(ctx)
		if err != nil {
			slog.Error("failed to count links for metric", "err", err)
			return math.NaN()
		}
		return float64(n)
	})
}
