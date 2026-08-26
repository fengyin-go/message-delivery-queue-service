package store

import "sync"

type ConsumerTagArchiveStore struct {
	mu      sync.RWMutex
	payload []byte
}

func NewConsumerTagArchiveStore() *ConsumerTagArchiveStore { return &ConsumerTagArchiveStore{} }

// Put 写入归档负载。store 持有独立副本，与调用方切片彻底解耦，
// 避免调用方复用原切片时改写已归档内容。
func (s *ConsumerTagArchiveStore) Put(payload []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.payload = append([]byte(nil), payload...)
}

// Snapshot 返回归档负载的副本，调用方修改返回值不影响 store 内部状态。
func (s *ConsumerTagArchiveStore) Snapshot() []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]byte(nil), s.payload...)
}
