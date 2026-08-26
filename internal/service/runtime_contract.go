package service

import (
	"errors"
	"fmt"
	"message-queue/internal/store"
)

type ProducerSendRetryCoordinator struct{ backend *store.ProducerSendRetryStore }

func NewProducerSendRetryCoordinator(b *store.ProducerSendRetryStore) *ProducerSendRetryCoordinator {
	return &ProducerSendRetryCoordinator{backend: b}
}
func (c *ProducerSendRetryCoordinator) Send(key string) error {
	err := c.backend.Attempt(key)
	if err == nil {
		return nil
	}
	// 用 %w 包裹，保留错误链，使临时故障可被 errors.As 识别并触发重试。
	var temporary *store.ProducerSendRetryTemporaryError
	if errors.As(err, &temporary) {
		return c.backend.Attempt(key)
	}
	return fmt.Errorf("dispatch failed: %w", err)
}
