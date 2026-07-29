package shortener

import (
	"context"
	"errors"
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

const maxSaveRetries = 5

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

	for i := 1; i <= maxSaveRetries; i++ {

		link, err := NewLink(originalUrl, s.codeLength)

		if err != nil {
			return nil, err
		}

		err = s.store.Save(ctx, link)

		if err != nil {
			if errors.Is(err, ErrUniqueConflict) {
				if i == maxSaveRetries {
					return nil, fmt.Errorf("save code error: %w", ErrUniqueConflict)
				}
				continue
			}
			return nil, err
		}
		return &link, nil
	}

	return nil, fmt.Errorf("save code error: %w", ErrUniqueConflict)
}

func (s *Service) GetLink(ctx context.Context, code string) (*Link, error) {

	link, err := s.store.Find(ctx, code)
	if err != nil {
		return nil, err
	}

	return link, nil
}
