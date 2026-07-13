package shortener

import (
	"errors"
	"testing"
)

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

func TestCreateLinkErrors(t *testing.T) {

	cases := []struct {
		name    string
		url     string
		length  int
		wantErr error
	}{
		{name: "invalid_length", url: "https://google.com/", length: 0, wantErr: ErrInvalidCodeLength},
		{name: "invalid_url", url: "not-a-url", length: 5, wantErr: ErrInvalidURL},
		{name: "empty_url", url: "", length: 4, wantErr: ErrInvalidURL},
		{name: "invalid_schema", url: "ftp://example.com", length: 5, wantErr: ErrInvalidURL},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			mem := NewMemoryStore()
			serv := NewService(mem, tc.length)

			_, err := serv.CreateLink(tc.url)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("expect err %v, got %v", tc.wantErr, err)
			}

			if len(mem.links) != 0 {
				t.Error("expect empty links in mem, got links with elements")
			}
		})
	}

}

func TestGetLinkValidCode(t *testing.T) {
	url := "https://google.com/"
	length := 15

	mem := NewMemoryStore()
	serv := NewService(mem, length)

	link, err := serv.CreateLink(url)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, ok := serv.GetLink(link.Code)

	if !ok {
		t.Fatal("already created link not found by code")
	}

	if *got != *link {
		t.Errorf("created and found links are different")
	}

}

func TestGetLinkInvalidCode(t *testing.T) {
	length := 15

	mem := NewMemoryStore()
	serv := NewService(mem, length)
	_, ok := serv.GetLink("missing")

	if ok {
		t.Fatal("expected not found with invalid code, got link")
	}
}
