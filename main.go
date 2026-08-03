package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
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

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("Config load", "error", err)
		return
	}
	pool, err := db.Connect(context.Background(), cfg.DSN)

	if err != nil {
		slog.Error("Database connect", "error", err)
		return
	}

	defer pool.Close()

	// store := shortener.NewMemoryStore()
	store := shortener.NewPostgresStore(pool)
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddress,
		Password: "",
		DB:       0,
	})
	defer rdb.Close()
	_, err = rdb.Ping(context.Background()).Result()

	if err != nil {
		slog.Error("RDB ping", "error", err)
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
		Handler: httpapi.RequestID(httpapi.Logging(mux)),
	}

	ch := make(chan error, 1)

	go func() {
		slog.Info("Running", "address", cfg.HTTPAddr)
		ch <- server.ListenAndServe()
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-ch:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Listen", "error", err)
		}
		return
	case <-ctx.Done():
		slog.Info("Shutdown signal received. Starting graceful termination...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("Shutdown", "error", err)
	}

	if err := <-ch; err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("Listen after shutdown", "error", err)
	}

	slog.Info("Exited cleanly.")

}
