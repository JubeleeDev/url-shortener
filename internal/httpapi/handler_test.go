package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/JubeleeDev/url-shortener/internal/shortener"
)

func TestHandlerSuccessCreateLink(t *testing.T) {
	wantURL := "https://example.com"
	body := strings.NewReader(fmt.Sprintf(`{"url":%q}`, wantURL))
	req := httptest.NewRequest(http.MethodPost, "/api/links", body)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	mux := initializeServeMux(8)
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expect status created, got %v", rr.Code)
	}

	var resp linkResponse
	err := json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	if resp.OriginalURL != wantURL {
		t.Errorf("expected URL %v, got %v", wantURL, resp.OriginalURL)
	}

	if resp.Code == "" {
		t.Error("got empty code in response")
	}

	if resp.Path != "/"+resp.Code {
		t.Errorf("expected path form \"/\" + %v, got \"/\" + %v ", resp.Code, resp.Path)
	}
}
func TestHandlerFailedCreateLink(t *testing.T) {
	t.Run("bad_json", func(t *testing.T) {
		body := strings.NewReader(`{"url":`)
		req := httptest.NewRequest(http.MethodPost, "/api/links", body)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		mux := initializeServeMux(8)
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status Bad request (not valid body), got %v", rr.Code)
		}
	})

	t.Run("zero_code_length", func(t *testing.T) {
		body := strings.NewReader(`{"url": "https:/google.com/"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/links", body)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		mux := initializeServeMux(0)
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status Bad request (code length = 0), got %v", rr.Code)
		}
	})

}
func TestHandlerSuccessGetLink(t *testing.T) {
	wantURL := "https://example.com"
	body := strings.NewReader(fmt.Sprintf(`{"url":%q}`, wantURL))
	req := httptest.NewRequest(http.MethodPost, "/api/links", body)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	mux := initializeServeMux(8)
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expect status created, got %v", rr.Code)
	}

	var resp linkResponse
	err := json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/links/"+resp.Code, nil)
	rr = httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expect status OK, got %v", rr.Code)
	}
	var getresp linkResponse
	err = json.NewDecoder(rr.Body).Decode(&getresp)
	if err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	if getresp.Code != resp.Code {
		t.Errorf("expected code %v, got %v", resp.Code, getresp.Code)
	}

	if getresp.OriginalURL != resp.OriginalURL {
		t.Errorf("expected URL %v, got %v", resp.OriginalURL, getresp.OriginalURL)
	}
	if getresp.Path != resp.Path {
		t.Errorf("expected path %v, got %v", resp.Path, getresp.Path)
	}

}
func TestHandlerGetLinkMissingCode(t *testing.T) {
	mux := initializeServeMux(8)

	req := httptest.NewRequest(http.MethodGet, "/api/links/missing", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expect status not found, got %v", rr.Code)
	}
}
func TestHandlerSuccessRedirect(t *testing.T) {
	wantURL := "https://example.com"
	body := strings.NewReader(fmt.Sprintf(`{"url":%q}`, wantURL))
	req := httptest.NewRequest(http.MethodPost, "/api/links", body)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	mux := initializeServeMux(8)
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expect status created, got %v", rr.Code)
	}

	var resp linkResponse
	err := json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	req = httptest.NewRequest(http.MethodGet, "/"+resp.Code, nil)
	rr = httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusFound {
		t.Fatalf("expected status found (302), got %v", rr.Code)
	}

	if loc := rr.Header().Get("Location"); loc != wantURL {
		t.Errorf("expected Location %v, got %v", wantURL, loc)
	}
}
func TestHandlerRedirectMissingLink(t *testing.T) {
	mux := initializeServeMux(8)

	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expect status not found, got %v", rr.Code)
	}
}

func initializeServeMux(length int) *http.ServeMux {

	store := shortener.NewMemoryStore()
	service := shortener.NewService(store, length)

	h := NewHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/links", h.CreateLink)
	mux.HandleFunc("GET /api/links/{code}", h.GetLink)
	mux.HandleFunc("GET /{code}", h.Redirect)

	return mux

}
