package store

import "sync"

type PayloadArchiveStore struct {
	mu      sync.RWMutex
	payload []byte
}

func NewPayloadArchiveStore() *PayloadArchiveStore { return &PayloadArchiveStore{} }
func (s *PayloadArchiveStore) Put(payload []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.payload = payload
}
func (s *PayloadArchiveStore) Snapshot() []byte { s.mu.RLock(); defer s.mu.RUnlock(); return s.payload }
