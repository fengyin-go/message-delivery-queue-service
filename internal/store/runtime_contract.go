package store

import (
	"context"
	"sync"
	"time"
)

type ConsumerSessionProbeStore struct {
	mu    sync.Mutex
	first context.Context
	delay time.Duration
}

func NewConsumerSessionProbeStore(delay time.Duration) *ConsumerSessionProbeStore {
	return &ConsumerSessionProbeStore{delay: delay}
}
func (s *ConsumerSessionProbeStore) Wait(ctx context.Context, key string) error {
	s.mu.Lock()
	if s.first == nil {
		s.first = ctx
	}
	active := s.first
	s.mu.Unlock()
	timer := time.NewTimer(s.delay)
	defer timer.Stop()
	select {
	case <-active.Done():
		return active.Err()
	case <-timer.C:
		return nil
	}
}
