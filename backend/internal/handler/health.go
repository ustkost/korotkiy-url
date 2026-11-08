package handler

import (
	"context"
	"net/http"
	"time"
)

type HealthHandler struct {
	isReady func() bool
	ping    func(ctx context.Context) error
}

func NewHealthHandler(isReady func() bool, ping func(ctx context.Context) error) *HealthHandler {
	return &HealthHandler{isReady: isReady, ping: ping}
}

func (h *HealthHandler) Live(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	if !h.isReady() {
		http.Error(w, "not ready", http.StatusServiceUnavailable)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()
	if err := h.ping(ctx); err != nil {
		http.Error(w, "db unavailable", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
}
