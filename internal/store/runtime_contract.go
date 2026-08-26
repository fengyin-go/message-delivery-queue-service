package store

import "sync"

type MessageRateSnapshotStore struct {
	mu    sync.Mutex
	value int
}

func NewMessageRateSnapshotStore(value int) *MessageRateSnapshotStore {
	return &MessageRateSnapshotStore{value: value}
}
func (s *MessageRateSnapshotStore) Snapshot() *int {
	s.mu.Lock()
	defer s.mu.Unlock()
	// 返回独立拷贝而非 &s.value，否则并发 Increment 会改写快照指向的字段。
	v := s.value
	return &v
}
func (s *MessageRateSnapshotStore) Increment()     { s.mu.Lock(); s.value++; s.mu.Unlock() }
