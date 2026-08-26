package store

import (
	"fmt"
	"sync"
)

type ProducerSendRetryTemporaryError struct{ Key string }

func (e *ProducerSendRetryTemporaryError) Error() string {
	return fmt.Sprintf("temporary delivery for %s", e.Key)
}

type ProducerSendRetryStore struct {
	mu      sync.Mutex
	calls   int
	records []string
}

func NewProducerSendRetryStore() *ProducerSendRetryStore { return &ProducerSendRetryStore{} }
func (s *ProducerSendRetryStore) Attempt(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls++
	if s.calls == 1 {
		// 首次返回临时故障，但不落最终发送记录，避免重试成功后留下重复记录。
		return &ProducerSendRetryTemporaryError{Key: key}
	}
	// 仅在发送成功时记录一次有效发送。
	s.records = append(s.records, key)
	return nil
}
func (s *ProducerSendRetryStore) Records() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.records...)
}
