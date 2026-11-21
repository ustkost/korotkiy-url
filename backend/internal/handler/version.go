package handler

import (
	"encoding/json"
	"net/http"
	"os"
)

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"version": os.Getenv("APP_VERSION"),
		"pod":     hostname,
	})
}
