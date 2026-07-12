package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/JubeleeDev/url-shortener/internal/shortener"
)

type createLinkRequest struct {
	URL string `json:"url"`
}

type createLinkResponse struct {
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
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req createLinkRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	link, err := h.service.CreateLink(req.URL)

	if err != nil {
		http.Error(w, "link was not created", http.StatusBadRequest)
		return
	}

	resp := createLinkResponse{
		OriginalURL: link.OriginalURL,
		Code:        link.Code,
		Path:        link.Path(), // метод у Link, который у тебя уже есть
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&resp)
}
