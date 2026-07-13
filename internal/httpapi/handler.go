package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/JubeleeDev/url-shortener/internal/shortener"
)

type createLinkRequest struct {
	URL string `json:"url"`
}

type linkResponse struct {
	OriginalURL string `json:"original_url"`
	Code        string `json:"code"`
	Path        string `json:"path"`
}

type Handler struct {
	service *shortener.Service
}

func NewHandler(service *shortener.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateLink(w http.ResponseWriter, r *http.Request) {

	var req createLinkRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	link, err := h.service.CreateLink(req.URL)

	if err != nil {
		if errors.Is(err, shortener.ErrInvalidURL) {
			http.Error(w, "invalid url", http.StatusBadRequest)
			return
		} else if errors.Is(err, shortener.ErrInvalidCodeLength) {
			http.Error(w, "invalid code length", http.StatusInternalServerError)
			return
		} else {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}

	resp := linkResponse{
		OriginalURL: link.OriginalURL,
		Code:        link.Code,
		Path:        link.Path(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&resp)
}

func (h *Handler) GetLink(w http.ResponseWriter, r *http.Request) {

	code := r.PathValue("code")

	if code == "" {
		http.Error(w, "code is empty", http.StatusBadRequest)
		return
	}

	link, ok := h.service.GetLink(code)

	if !ok {
		http.Error(w, "link not found", http.StatusNotFound)
		return
	}

	resp := linkResponse{
		OriginalURL: link.OriginalURL,
		Code:        link.Code,
		Path:        link.Path(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	if code == "" {
		http.Error(w, "code is empty", http.StatusBadRequest)
		return
	}

	link, ok := h.service.GetLink(code)

	if !ok {
		http.Error(w, "link not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, link.OriginalURL, http.StatusFound)
}
