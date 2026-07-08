package shortener

import "testing"

func TestNewLinkCreatesLink(t *testing.T) {
	originalURL := "https://google.com/"
	length := 19

	link, err := NewLink(originalURL, length)

	if err != nil {
		t.Fatalf("got unexpected error: %v", err)
	}

	if link.OriginalURL != originalURL {
		t.Errorf("expected original url %v, got %v", originalURL, link.OriginalURL)
	}

	if len(link.Code) != length {
		t.Errorf("expected code length %d, got %d", length, len(link.Code))
	}
}

func TestNewLinkReturnsErrorForInvalidLength(t *testing.T) {
	originalURL := "https://google.com/"
	length := 0

	_, err := NewLink(originalURL, length)

	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestLinkPathReturnsCodeWithLeadingSlash(t *testing.T) {

	link := Link{
		Code: "abc123",
	}

	if link.Path() != "/"+link.Code {
		t.Errorf("expected link path /%v, got %v", link.Code, link.Path())
	}

}
