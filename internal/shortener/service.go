package shortener

import (
	"context"
	"fmt"
	"net/url"
)

type Service struct {
	store      Store
	codeLength int
}

type Store interface {
	Save(ctx context.Context, link Link) error
	Find(ctx context.Context, code string) (*Link, error)
}

func NewService(store Store, codeLen int) *Service {
	return &Service{store: store, codeLength: codeLen}
}

func (s *Service) CreateLink(ctx context.Context, originalUrl string) (*Link, error) {

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

	err = s.store.Save(ctx, link)

	if err != nil {
		return nil, err
	}

	return &link, nil

}

func (s *Service) GetLink(ctx context.Context, code string) (*Link, error) {

	link, err := s.store.Find(ctx, code)
	if err != nil {
		return nil, err
	}

	return link, nil
}
