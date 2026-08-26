package store

import "sync"

type MessageRateSnapshotStore struct {
	mu    sync.Mutex
	value int
}

func NewMessageRateSnapshotStore(value int) *MessageRateSnapshotStore {
	return &MessageRateSnapshotStore{value: value}
}
func (s *MessageRateSnapshotStore) Snapshot() *int { s.mu.Lock(); defer s.mu.Unlock(); return &s.value }
func (s *MessageRateSnapshotStore) Increment()     { s.mu.Lock(); s.value++; s.mu.Unlock() }
