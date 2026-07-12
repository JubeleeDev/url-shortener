package shortener

import "testing"

func TestSuccessCreateLink(t *testing.T) {
	url := "https://google.com/"
	length := 15

	mem := NewMemoryStore()
	serv := NewService(mem, length)

	link, err := serv.CreateLink(url)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if link.OriginalURL != url {
		t.Errorf("expected url %v, got %v", url, link.OriginalURL)
	}

	if link.Code == "" {
		t.Errorf("expected not empty code")
	}

	if len(link.Code) != length {
		t.Errorf("expected len code %d, got %d", length, len(link.Code))
	}

	got, ok := mem.Find(link.Code)
	if !ok {
		t.Errorf("expected link with code %v, got not found", link.Code)
	}

	if *got != *link {
		t.Errorf("created and found links are different")
	}
}

func TestCreateLinkInvalidLength(t *testing.T) {
	url := "https://google.com/"
	length := 0

	mem := NewMemoryStore()
	serv := NewService(mem, length)

	_, err := serv.CreateLink(url)

	if err == nil {
		t.Fatalf("expect error, got nil")
	}

	if len(mem.links) != 0 {
		t.Error("expect empty links in mem, got links with elements")
	}
}
