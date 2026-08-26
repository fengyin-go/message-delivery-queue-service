package store

import (
	"context"
	"sync"
	"time"
)

type SubscriptionSyncCancelStore struct {
	mu    sync.Mutex
	calls int
	delay time.Duration
}

func NewSubscriptionSyncCancelStore(delay time.Duration) *SubscriptionSyncCancelStore {
	return &SubscriptionSyncCancelStore{delay: delay}
}
func (s *SubscriptionSyncCancelStore) Attempt(ctx context.Context) error {
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
func (s *SubscriptionSyncCancelStore) Calls() int { s.mu.Lock(); defer s.mu.Unlock(); return s.calls }
