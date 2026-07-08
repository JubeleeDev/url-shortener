package shortener

import "testing"

func TestCreateLinkWithValidCredentials(t *testing.T) {
	originalURL := "https://google.com/"
	length := 19

	link, err := NewLink(originalURL, length)

	if err != nil {
		t.Fatalf("got unexpected error: %v", err)
	}

	if link.OriginalURL != originalURL {
		t.Errorf("original url is %v, after created code got %v", originalURL, link.OriginalURL)
	}

	if len(link.Code) != length {
		t.Errorf("original length is %d, got %d", length, len(link.Code))
	}
}

func TestCreateLinkWithInvalidCredentials(t *testing.T) {
	originalURL := "https://google.com/"
	length := 0

	_, err := NewLink(originalURL, length)

	if err == nil {
		t.Errorf("expected error, got nil")
	}

}
