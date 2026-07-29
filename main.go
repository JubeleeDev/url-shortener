package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

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
	defer rdb.Close()
	_, err = rdb.Ping(context.Background()).Result()

	if err != nil {
		fmt.Println(err)
		return
	}

	linkCache := cache.NewCache(rdb)
	service := shortener.NewService(store, cfg.CodeLength, linkCache)

	h := httpapi.NewHandler(service)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/links", h.CreateLink)
	mux.HandleFunc("GET /api/links/{code}", h.GetLink)
	mux.HandleFunc("GET /{code}", h.Redirect)

	server := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: mux,
	}

	ch := make(chan error, 1)

	go func() {
		fmt.Println("server is running at", cfg.HTTPAddr)
		ch <- server.ListenAndServe()
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-ch:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Server listen error: %v\n", err)
		}
		return
	case <-ctx.Done():
		log.Println("Shutdown signal received. Starting graceful termination...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v\n", err)
	}

	if err := <-ch; err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("Server listen error: %v", err)
	}

	log.Println("Server exited cleanly.")

}
