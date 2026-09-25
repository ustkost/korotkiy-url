package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ustkost/korotkiy-url/internal/service"
)

type ClickHandler struct {
	service *service.ClickService
}

func NewClickHandler(s *service.ClickService) *ClickHandler {
	return &ClickHandler{service: s}
}

func (h *ClickHandler) List(w http.ResponseWriter, r *http.Request) {
	id, err := parseInt64(r.PathValue("id"))
	if err != nil {
		http.Error(w, "id must be a valid integer", http.StatusBadRequest)
		return
	}

	limit, err := parseIntParam(r, "limit", defaultLimit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	offset, err := parseIntParam(r, "offset", 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	clicks, err := h.service.ListByLinkID(r.Context(), id, limit, offset)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clicks)
}
