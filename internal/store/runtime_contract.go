package store

import (
	"fmt"
	"sync"
)

type RetryDispatchTemporaryError struct{ Key string }

func (e *RetryDispatchTemporaryError) Error() string {
	return fmt.Sprintf("temporary delivery for %s", e.Key)
}

type RetryDispatchStore struct {
	mu      sync.Mutex
	calls   int
	records []string
}

func NewRetryDispatchStore() *RetryDispatchStore { return &RetryDispatchStore{} }
func (s *RetryDispatchStore) Attempt(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls++
	s.records = append(s.records, key)
	if s.calls == 1 {
		return &RetryDispatchTemporaryError{Key: key}
	}
	s.records = append(s.records, key)
	return nil
}
func (s *RetryDispatchStore) Records() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.records...)
}
