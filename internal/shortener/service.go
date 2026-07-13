package shortener

import (
	"fmt"
	"net/url"
)

type Service struct {
	store      Store
	codeLength int
}

type Store interface {
	Save(link Link)
	Find(code string) (*Link, bool)
}

func NewService(store Store, codeLen int) *Service {
	return &Service{store: store, codeLength: codeLen}
}

func (s *Service) CreateLink(originalUrl string) (*Link, error) {

	u, err := url.ParseRequestURI(originalUrl)

	if err != nil {
		return nil, fmt.Errorf("error: %w", ErrInvalidURL)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, ErrInvalidURL
	}
	if u.Host == "" {
		return nil, ErrInvalidURL
	}

	link, err := NewLink(originalUrl, s.codeLength)

	if err != nil {
		return nil, err
	}

	s.store.Save(link)

	return &link, nil

}

func (s *Service) GetLink(code string) (*Link, bool) {
	if len(code) == 0 {
		return nil, false
	}

	link, ok := s.store.Find(code)
	if !ok {
		return nil, false
	}

	return link, true
}
