package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/ustkost/korotkiy-url/internal/metrics"
	"github.com/ustkost/korotkiy-url/internal/repository"
	"github.com/ustkost/korotkiy-url/internal/service"
)

type RedirectHandler struct {
	linkService  *service.LinkService
	clickService *service.ClickService
}

func NewRedirectHandler(linkService *service.LinkService, clickService *service.ClickService) *RedirectHandler {
	return &RedirectHandler{
		linkService:  linkService,
		clickService: clickService,
	}
}

func (h *RedirectHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		metrics.RedirectDuration.Observe(time.Since(start).Seconds())
	}()
	shortCode := r.PathValue("shortCode")

	link, err := h.linkService.GetByShortCode(r.Context(), shortCode)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			metrics.RedirectsTotal.WithLabelValues("not_found").Inc()
		} else {
			metrics.RedirectsTotal.WithLabelValues("error").Inc()
		}
		writeServiceError(w, err)
		return
	}

	referrer := r.Header.Get("Referer")
	go func() {
		if err := h.clickService.RecordClick(context.Background(), link.ID, referrer); err != nil {
			slog.Error("failed to record click for link", "link_id", link.ID, "err", err)
		}
	}()

	http.Redirect(w, r, link.OriginalURL, http.StatusFound)
	metrics.RedirectsTotal.WithLabelValues("found").Inc()
}
