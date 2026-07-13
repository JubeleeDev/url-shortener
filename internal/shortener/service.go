package shortener

import (
	"fmt"
	"net/url"
)

type Service struct {
	memory     *MemoryStore
	codeLength int
}

func NewService(mem *MemoryStore, codeLen int) *Service {
	return &Service{memory: mem, codeLength: codeLen}
}

func (s *Service) CreateLink(originalUrl string) (*Link, error) {

	u, err := url.ParseRequestURI(originalUrl)

	if err != nil {
		return nil, fmt.Errorf("invalid url: %w", err)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("invalid request scheme")
	}
	if u.Host == "" {
		return nil, fmt.Errorf("host is empty")
	}

	link, err := NewLink(originalUrl, s.codeLength)

	if err != nil {
		return nil, err
	}

	s.memory.Save(link)

	return &link, nil

}

func (s *Service) GetLink(code string) (*Link, bool) {
	if len(code) == 0 {
		return nil, false
	}

	link, ok := s.memory.Find(code)
	if !ok {
		return nil, false
	}

	return link, true
}
