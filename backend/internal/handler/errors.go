package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/ustkost/korotkiy-url/internal/repository"
	"github.com/ustkost/korotkiy-url/internal/service"
)

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidURL),
		errors.Is(err, service.ErrInvalidShortCode),
		errors.Is(err, service.ErrInvalidLimit),
		errors.Is(err, service.ErrInvalidOffset):

		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, repository.ErrDuplicateCode):
		http.Error(w, err.Error(), http.StatusConflict)
	case errors.Is(err, repository.ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, service.ErrCodeGenerationFailed):
		http.Error(w, "unable to generate a unique short code, please try again", http.StatusInternalServerError)
	default:
		log.Printf("unexpected error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
