package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ustkost/korotkiy-url/internal/service"
)

type LinkHandler struct {
	service *service.LinkService
}

func NewLinkHandler(s *service.LinkService) *LinkHandler {
	return &LinkHandler{service: s}
}

type createRequest struct {
	OriginalURL string `json:"original_url"`
	CustomCode  string `json:"custom_code"`
}

func (h *LinkHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	link, err := h.service.Create(r.Context(), req.OriginalURL, req.CustomCode)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(link)
}

func (h *LinkHandler) List(w http.ResponseWriter, r *http.Request) {
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

	links, err := h.service.List(r.Context(), limit, offset)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(links)
}

func (h *LinkHandler) Get(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("shortCode")

	link, err := h.service.GetByShortCode(r.Context(), shortCode)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(link)
}

type updateOriginalURLRequest struct {
	OriginalURL string `json:"original_url"`
}

func (h *LinkHandler) UpdateOriginalURL(w http.ResponseWriter, r *http.Request) {
	id, err := parseInt64(r.PathValue("id"))
	if err != nil {
		http.Error(w, "id must be a valid integer", http.StatusBadRequest)
		return
	}

	var req updateOriginalURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	link, err := h.service.UpdateOriginalURL(r.Context(), id, req.OriginalURL)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(link)
}

type updateShortCodeRequest struct {
	CustomCode string `json:"custom_code"`
}

func (h *LinkHandler) UpdateShortCode(w http.ResponseWriter, r *http.Request) {
	id, err := parseInt64(r.PathValue("id"))
	if err != nil {
		http.Error(w, "id must be a valid integer", http.StatusBadRequest)
		return
	}

	var req updateShortCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	link, err := h.service.UpdateShortCode(r.Context(), id, req.CustomCode)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(link)
}

func (h *LinkHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseInt64(r.PathValue("id"))
	if err != nil {
		http.Error(w, "id must be a valid integer", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		writeServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
