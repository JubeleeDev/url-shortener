package httpapi

import (
	"net/http"
	"testing"

	"github.com/JubeleeDev/url-shortener/internal/httpapi"
	"github.com/JubeleeDev/url-shortener/internal/shortener"
)

func TestHandlerSuccessCreateLink(t *testing.T) {

}
func TestHandlerUnsuccessCreateLink(t *testing.T) {

}
func TestHandlerSuccessGetLink(t *testing.T) {

}
func TestHandlerGetLinkMissingCode(t *testing.T) {

}
func TestHandlerSuccessRedirect(t *testing.T) {

}
func TestHandlerRedirectMissingLink(t *testing.T) {

}

func initializeServeMux() *http.ServeMux {

	store := shortener.NewMemoryStore()
	service := shortener.NewService(store, 8)

	h := httpapi.NewHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/links", h.CreateLink)
	mux.HandleFunc("GET /api/links/{code}", h.GetLink)
	mux.HandleFunc("GET /{code}", h.Redirect)

}
