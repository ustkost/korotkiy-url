package main

import (
	"log"
	"net/http"

	"github.com/ustkost/korotkiy-url/internal/config"
	"github.com/ustkost/korotkiy-url/internal/handler"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	mux := handler.NewRouter()

	log.Printf("listening on :%s\n", cfg.Port);
	log.Fatal(http.ListenAndServe(":" + cfg.Port, mux))
}
