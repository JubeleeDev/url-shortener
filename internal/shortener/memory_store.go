package shortener

type MemoryStore struct {
	links map[string]Link
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{links: make(map[string]Link)}
}

func (s *MemoryStore) Save(link Link) {
	s.links[link.Code] = link
}

func (s *MemoryStore) Find(code string) (Link, bool) {
	value, ok := s.links[code]

	if !ok {
		return Link{}, false
	}

	return value, true
}
