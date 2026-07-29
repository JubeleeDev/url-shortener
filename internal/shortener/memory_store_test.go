package shortener

import (
	"testing"
)

func TestMemoryStoreSaveAndFind(t *testing.T) {
	store := NewMemoryStore()
	link := Link{
		Code:        "abc123",
		OriginalURL: "https://example.com",
	}
	err := store.Save(t.Context(), link)

	if err != nil {
		t.Fatalf("got an unexpected error %v", err)
	}

	foundLink, err := store.Find(t.Context(), link.Code)
	if err != nil {
		t.Fatalf("expected link to be found, got %v", err)
	}

	if link.Code != foundLink.Code {
		t.Errorf("expected code %v, got %v", link.Code, foundLink.Code)
	}

	if link.OriginalURL != foundLink.OriginalURL {
		t.Errorf("expected original URL %v, got %v", link.OriginalURL, foundLink.OriginalURL)
	}
}

func TestMemoryStoreFindLinkByUnknownCode(t *testing.T) {
	store := NewMemoryStore()

	_, err := store.Find(t.Context(), "missing")

	if err == nil {
		t.Error("expected error, got link")
	}
}
