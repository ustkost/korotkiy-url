package handler

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/ustkost/korotkiy-url/internal/metrics"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Unwrap() http.ResponseWriter {
	return r.ResponseWriter
}

const (
	PathLive  = "/livez"
	PathReady = "/readyz"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r)

		isProbe := r.URL.Path == PathLive || r.URL.Path == PathReady
		if isProbe && rec.status == http.StatusOK {
			return
		}

		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}

func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == PathLive || r.URL.Path == PathReady {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()
		metrics.HttpInFlight.Inc()
		defer metrics.HttpInFlight.Dec()

		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)

		route := r.Pattern
		if route == "" {
			route = "unmatched"
		}

		metrics.HttpRequests.WithLabelValues(r.Method, route, strconv.Itoa(rec.status)).Inc()
		metrics.HttpDuration.WithLabelValues(r.Method, route).Observe(time.Since(start).Seconds())
	})
}
