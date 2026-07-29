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
	cache      Cache
}

type Store interface {
	Save(ctx context.Context, link Link) error
	Find(ctx context.Context, code string) (*Link, error)
}

type Cache interface {
	Get(ctx context.Context, code string) (string, error)
	Set(ctx context.Context, code string, url string, ttl int) error
	Delete(ctx context.Context, code string) error
}

func NewService(store Store, codeLen int, cache Cache) *Service {
	return &Service{store: store, codeLength: codeLen, cache: cache}
}

const maxSaveRetries = 5
const ttlSeconds = 3600

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
		_ = s.cache.Set(ctx, link.Code, link.OriginalURL, ttlSeconds)

		return &link, nil
	}

	return nil, fmt.Errorf("save code error: %w", ErrUniqueConflict)
}

func (s *Service) GetLink(ctx context.Context, code string) (*Link, error) {

	url, err := s.cache.Get(ctx, code)

	if err == nil {
		return &Link{Code: code, OriginalURL: url}, nil
	}

	link, err := s.store.Find(ctx, code)
	if err != nil {
		return nil, err
	}

	_ = s.cache.Set(ctx, code, link.OriginalURL, ttlSeconds)

	return link, nil
}
