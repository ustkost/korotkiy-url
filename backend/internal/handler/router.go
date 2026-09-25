package handler

import "net/http"

func NewRouter(linkHandler *LinkHandler, clickHandler *ClickHandler, redirectHandler *RedirectHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", Health)

	mux.HandleFunc("POST /links", linkHandler.Create)
	mux.HandleFunc("GET /links", linkHandler.List)
	mux.HandleFunc("GET /links/{shortCode}", linkHandler.Get)
	mux.HandleFunc("PATCH /links/{id}/url", linkHandler.UpdateOriginalURL)
	mux.HandleFunc("PATCH /links/{id}/code", linkHandler.UpdateShortCode)
	mux.HandleFunc("DELETE /links/{id}", linkHandler.Delete)

	mux.HandleFunc("GET /links/{id}/clicks", clickHandler.List)

	mux.HandleFunc("GET /s/{shortCode}", redirectHandler.Redirect)

	return mux
}
