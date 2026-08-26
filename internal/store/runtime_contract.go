package store

import (
	"context"
	"sync"
	"time"
)

type ProducerCancellationStore struct {
	mu    sync.Mutex
	calls int
	delay time.Duration
}

func NewProducerCancellationStore(delay time.Duration) *ProducerCancellationStore {
	return &ProducerCancellationStore{delay: delay}
}
func (s *ProducerCancellationStore) Attempt(ctx context.Context) error {
	s.mu.Lock()
	s.calls++
	s.mu.Unlock()
	timer := time.NewTimer(s.delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
func (s *ProducerCancellationStore) Calls() int { s.mu.Lock(); defer s.mu.Unlock(); return s.calls }
