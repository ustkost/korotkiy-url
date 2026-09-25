package main

import (
	"context"
	"log"
	"net/http"

	"github.com/ustkost/korotkiy-url/internal/config"
	"github.com/ustkost/korotkiy-url/internal/db"
	"github.com/ustkost/korotkiy-url/internal/handler"
	"github.com/ustkost/korotkiy-url/internal/repository"
	"github.com/ustkost/korotkiy-url/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	pool, err := db.Connect(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	linkRepo := repository.NewLinkRepository(pool)
	clickRepo := repository.NewClickRepository(pool)

	linkService := service.NewLinkService(linkRepo)
	clickService := service.NewClickService(clickRepo)

	linkHandler := handler.NewLinkHandler(linkService)
	clickHandler := handler.NewClickHandler(clickService)
	redirectHandler := handler.NewRedirectHandler(linkService, clickService)

	mux := handler.NewRouter(linkHandler, clickHandler, redirectHandler)

	log.Printf("listening on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}
