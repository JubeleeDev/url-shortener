package shortener

import "testing"

func TestMemoryStoreSaveAndFind(t *testing.T) {
	store := NewMemoryStore()
	link := Link{
		Code:        "abc123",
		OriginalURL: "https://example.com",
	}
	store.Save(link)

	foundedLink, ok := store.Find(link.Code)
	if !ok {
		t.Fatal("expected founded link, got empty")
	}

	if link.Code != foundedLink.Code {
		t.Errorf("expected founded link code == abc123, got %v", foundedLink.Code)
	}

	if link.OriginalURL != foundedLink.OriginalURL {
		t.Errorf("expected founded link original URL == https://example.com, got %v", foundedLink.OriginalURL)
	}
}

func TestMemoryStoreFindReturnsFalseForUnknownCode(t *testing.T) {
	store := NewMemoryStore()

	emptyLink := Link{}

	_, ok := store.Find(emptyLink.Code)

	if ok {
		t.Error("expected error, got founded link")
	}
}
