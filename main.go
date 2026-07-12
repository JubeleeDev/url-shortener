package main

import (
	"fmt"
	"net/http"

	"github.com/JubeleeDev/url-shortener/internal/httpapi"
	"github.com/JubeleeDev/url-shortener/internal/shortener"
)

func main() {
	store := shortener.NewMemoryStore()
	service := shortener.NewService(store, 8)

	h := httpapi.NewHandler(service)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/links", h.CreateLink)
	mux.HandleFunc("GET /api/links/{code}", h.GetLink)

	fmt.Println("server is running at 8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("unexpected error:", err)
	}
}
