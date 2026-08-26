package store

import "sync"

type TopicCounterSnapshotStore struct {
	mu    sync.Mutex
	value int
}

func NewTopicCounterSnapshotStore(value int) *TopicCounterSnapshotStore {
	return &TopicCounterSnapshotStore{value: value}
}
func (s *TopicCounterSnapshotStore) Snapshot() *int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return &s.value
}
func (s *TopicCounterSnapshotStore) Increment() { s.mu.Lock(); s.value++; s.mu.Unlock() }
