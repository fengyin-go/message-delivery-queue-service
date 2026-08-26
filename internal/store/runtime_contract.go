package store

import (
	"fmt"
	"sync"
)

type DeadLetterReplayTemporaryError struct{ Key string }

func (e *DeadLetterReplayTemporaryError) Error() string {
	return fmt.Sprintf("temporary delivery for %s", e.Key)
}

type DeadLetterReplayStore struct {
	mu      sync.Mutex
	calls   int
	records []string
}

func NewDeadLetterReplayStore() *DeadLetterReplayStore { return &DeadLetterReplayStore{} }
func (s *DeadLetterReplayStore) Attempt(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls++
	s.records = append(s.records, key)
	if s.calls == 1 {
		return &DeadLetterReplayTemporaryError{Key: key}
	}
	s.records = append(s.records, key)
	return nil
}
func (s *DeadLetterReplayStore) Records() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.records...)
}
