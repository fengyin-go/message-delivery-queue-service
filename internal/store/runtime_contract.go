package store

import (
	"fmt"
	"sync"
)

type ProducerSendRetryTemporaryError struct{ Key string }

func (e *ProducerSendRetryTemporaryError) Error() string {
	return fmt.Sprintf("temporary delivery for %s", e.Key)
}

type ProducerSendRetryStore struct {
	mu      sync.Mutex
	calls   int
	records []string
}

func NewProducerSendRetryStore() *ProducerSendRetryStore { return &ProducerSendRetryStore{} }
func (s *ProducerSendRetryStore) Attempt(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls++
	s.records = append(s.records, key)
	if s.calls == 1 {
		return &ProducerSendRetryTemporaryError{Key: key}
	}
	s.records = append(s.records, key)
	return nil
}
func (s *ProducerSendRetryStore) Records() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.records...)
}
