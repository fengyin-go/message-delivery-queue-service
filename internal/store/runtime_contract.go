package store

import (
	"context"
	"sync"
	"time"
)

type DeadLetterReplayCancelStore struct {
	mu    sync.Mutex
	calls int
	delay time.Duration
}

func NewDeadLetterReplayCancelStore(delay time.Duration) *DeadLetterReplayCancelStore {
	return &DeadLetterReplayCancelStore{delay: delay}
}
func (s *DeadLetterReplayCancelStore) Attempt(ctx context.Context) error {
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
func (s *DeadLetterReplayCancelStore) Calls() int { s.mu.Lock(); defer s.mu.Unlock(); return s.calls }
