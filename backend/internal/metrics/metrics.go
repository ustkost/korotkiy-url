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
