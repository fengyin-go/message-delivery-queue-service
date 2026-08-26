package store

import "sync"

type GroupOffsetSnapshotStore struct {
	mu    sync.Mutex
	value int
}

func NewGroupOffsetSnapshotStore(value int) *GroupOffsetSnapshotStore {
	return &GroupOffsetSnapshotStore{value: value}
}
func (s *GroupOffsetSnapshotStore) Snapshot() *int { s.mu.Lock(); defer s.mu.Unlock(); v := s.value; return &v }
func (s *GroupOffsetSnapshotStore) Increment()     { s.mu.Lock(); s.value++; s.mu.Unlock() }
