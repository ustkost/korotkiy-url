package handler

import (
	"context"
	"log"
	"net/http"

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
	shortCode := r.PathValue("shortCode")

	link, err := h.linkService.GetByShortCode(r.Context(), shortCode)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	referrer := r.Header.Get("Referer")
	go func() {
		if err := h.clickService.RecordClick(context.Background(), link.ID, referrer); err != nil {
			log.Printf("failed to record click for link %d: %v", link.ID, err)
		}
	}()

	http.Redirect(w, r, link.OriginalURL, http.StatusFound)
}

