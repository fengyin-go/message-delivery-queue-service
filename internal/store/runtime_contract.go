package store

import "sync"

type RetryReasonExportStore struct {
	mu      sync.RWMutex
	payload []byte
}

func NewRetryReasonExportStore() *RetryReasonExportStore { return &RetryReasonExportStore{} }
func (s *RetryReasonExportStore) Put(payload []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.payload = payload
}
func (s *RetryReasonExportStore) Snapshot() []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.payload
}
