package store

import "sync"

type ConsumerTagArchiveStore struct {
	mu      sync.RWMutex
	payload []byte
}

func NewConsumerTagArchiveStore() *ConsumerTagArchiveStore { return &ConsumerTagArchiveStore{} }
func (s *ConsumerTagArchiveStore) Put(payload []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.payload = payload
}
func (s *ConsumerTagArchiveStore) Snapshot() []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.payload
}
