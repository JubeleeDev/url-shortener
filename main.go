package main

import (
	"net/http"

	"github.com/JubeleeDev/url-shortener/internal/httpapi"
	"github.com/JubeleeDev/url-shortener/internal/shortener"
)

func main() {
	store := shortener.NewMemoryStore()
	service := shortener.NewService(store, 8)
	h := httpapi.NewHandler(service)
	mux := http.NewServeMux()
	mux.HandleFunc("/api/links", h.CreateLink)
	http.ListenAndServe(":8080", mux)
}
