package service

import (
	"log"
	"net/http"
	"newscrapper/handler"
	"newscrapper/internal/config"
)

func RunHTTPServer(di *config.DI) {
	h := handler.Handler{DI: *di}
	mux := http.NewServeMux()

	h.Mount(mux)

	server := http.Server{
		Addr:    "localhost:8000",
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}
