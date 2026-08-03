package httpapi

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"regexp"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

type ctxKey int

const requestIDKey ctxKey = 1

var requestIDPattern = regexp.MustCompile(`^[A-Za-z0-9._-]{1,128}$`)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w}
		next.ServeHTTP(rec, r)
		status := rec.status
		if status == 0 {
			status = http.StatusOK
		}

		duration := time.Since(start)
		slog.Info("request",
			"request_id", requestIDFromContext(r.Context()),
			"method", r.Method,
			"path", r.URL.Path,
			"status", status,
			"duration_ms", duration.Milliseconds())
	})
}

func withRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

func requestIDFromContext(ctx context.Context) string {
	id, ok := ctx.Value(requestIDKey).(string)
	if !ok {
		return ""
	}
	return id
}

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if !requestIDPattern.MatchString(id) {
			var err error
			id, err = newRequestID()
			if err != nil {
				slog.Error("generate request id", "error", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("X-Request-ID", id)
		ctx := withRequestID(r.Context(), id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func newRequestID() (string, error) {
	b := make([]byte, 16)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil

}
