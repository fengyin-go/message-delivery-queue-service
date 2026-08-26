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
	wrapped := fmt.Errorf("dispatch failed: %v", err)
	var temporary *store.ProducerSendRetryTemporaryError
	if errors.As(wrapped, &temporary) {
		return c.backend.Attempt(key)
	}
	return wrapped
}
