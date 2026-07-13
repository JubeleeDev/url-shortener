package main

import (
	"fmt"
	"net/http"

	"github.com/JubeleeDev/url-shortener/internal/config"
	"github.com/JubeleeDev/url-shortener/internal/httpapi"
	"github.com/JubeleeDev/url-shortener/internal/shortener"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println(err)
		return
	}
	store := shortener.NewMemoryStore()
	service := shortener.NewService(store, cfg.CodeLength)

	h := httpapi.NewHandler(service)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/links", h.CreateLink)
	mux.HandleFunc("GET /api/links/{code}", h.GetLink)
	mux.HandleFunc("GET /{code}", h.Redirect)

	fmt.Println("server is running at", cfg.HTTPAddr)
	err = http.ListenAndServe(cfg.HTTPAddr, mux)
	if err != nil {
		fmt.Println("unexpected error:", err)
	}
}
