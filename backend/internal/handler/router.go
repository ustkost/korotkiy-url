package handler

import "net/http"

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", Health)

	return mux
}
