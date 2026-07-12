package shortener

type Service struct {
	memory     *MemoryStore
	codeLength int
}

func NewService(mem *MemoryStore, codeLen int) *Service {
	return &Service{memory: mem, codeLength: codeLen}
}

func (s *Service) CreateLink(originalUrl string) (*Link, error) {

	link, err := NewLink(originalUrl, s.codeLength)

	if err != nil {
		return nil, err
	}

	s.memory.Save(link)

	return &link, nil

}
