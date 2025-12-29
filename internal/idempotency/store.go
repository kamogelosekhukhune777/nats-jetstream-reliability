package idempotency

import "sync"

// Simple in-memory store (replace with Redis in prod)
type Store struct {
	mu   sync.Mutex
	seen map[string]struct{}
}

func New() *Store {
	return &Store{seen: make(map[string]struct{})}
}

// 4. Idempotency keys
func (s *Store) Seen(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.seen[key]; ok {
		return true
	}

	s.seen[key] = struct{}{}
	return false
}
