package shortener

import (
	"context"
	"sync"
)

type MemoryStore struct {
	links map[string]Link
	mu    sync.Mutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{links: make(map[string]Link)}
}

func (s *MemoryStore) Save(ctx context.Context, link Link) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.links[link.Code]; ok {
		return ErrUniqueConflict
	}

	s.links[link.Code] = link
	return nil
}

func (s *MemoryStore) Find(ctx context.Context, code string) (*Link, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	value, ok := s.links[code]

	if !ok {
		return nil, ErrNotFound
	}
	return &value, nil
}
