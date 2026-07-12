package shortener

import "testing"

func TestMemoryStoreSaveAndFind(t *testing.T) {
	store := NewMemoryStore()
	link := Link{
		Code:        "abc123",
		OriginalURL: "https://example.com",
	}
	store.Save(link)

	foundLink, ok := store.Find(link.Code)
	if !ok {
		t.Fatal("expected link to be found")
	}

	if link.Code != foundLink.Code {
		t.Errorf("expected code %v, got %v", link.Code, foundLink.Code)
	}

	if link.OriginalURL != foundLink.OriginalURL {
		t.Errorf("expected original URL %v, got %v", link.OriginalURL, foundLink.OriginalURL)
	}
}

func TestMemoryStoreFindReturnsFalseForUnknownCode(t *testing.T) {
	store := NewMemoryStore()

	_, ok := store.Find("missing")

	if ok {
		t.Error("expected false, got true")
	}
}
