package idempotency

import "sync"

// Simple in-memory store
type Store struct {
	mu   sync.Mutex
	seen map[string]struct{}
}

func New() *Store {
	return &Store{
		seen: make(map[string]struct{}),
	}
}

func (s *Store) Seen(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.seen[key]
	if !ok {
		s.seen[key] = struct{}{}
	}
	return ok
}
