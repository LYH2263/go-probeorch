package result

// Store 每目标一个环形历史。
type Store struct {
	cap   int
	rings map[string]*Ring
}

func NewStore(capacity int) *Store {
	return &Store{cap: capacity, rings: make(map[string]*Ring)}
}

func (s *Store) Ensure(id string) {
	if _, ok := s.rings[id]; !ok {
		s.rings[id] = NewRing(s.cap)
	}
}

func (s *Store) Push(id string, e Entry) {
	s.Ensure(id)
	s.rings[id].Push(e)
}

func (s *Store) List(id string) []Entry {
	r, ok := s.rings[id]
	if !ok {
		return nil
	}
	return r.List()
}

func (s *Store) Drop(id string) {
	delete(s.rings, id)
}

func (s *Store) Clear() {
	s.rings = make(map[string]*Ring)
}
