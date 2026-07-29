package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/JubeleeDev/url-shortener/db"
	"github.com/JubeleeDev/url-shortener/internal/cache"
	"github.com/JubeleeDev/url-shortener/internal/config"
	"github.com/JubeleeDev/url-shortener/internal/httpapi"
	"github.com/JubeleeDev/url-shortener/internal/shortener"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println(err)
		return
	}
	pool, err := db.Connect(context.Background(), cfg.DSN)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer pool.Close()

	// store := shortener.NewMemoryStore()
	store := shortener.NewPostgresStore(pool)
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddress, // Redis server address
		Password: "",               // No password by default
		DB:       0,                // Use default database ID 0
	})
	cache := cache.NewCache(rdb)
	service := shortener.NewService(store, cfg.CodeLength, cache)

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
