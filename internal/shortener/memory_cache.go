package shortener

import (
	"context"
)

type MemoryCache struct{}

func (m *MemoryCache) Get(ctx context.Context, code string) (string, error) {
	return "", ErrCacheMiss
}

func (m *MemoryCache) Set(ctx context.Context, code string, url string, ttl int) error {
	return nil
}

func (m *MemoryCache) Delete(ctx context.Context, code string) error {
	return nil
}
