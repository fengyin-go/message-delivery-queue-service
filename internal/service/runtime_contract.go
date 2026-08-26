package service

import (
	"errors"
	"fmt"
	"message-queue/internal/store"
)

type RetryDispatchCoordinator struct{ backend *store.RetryDispatchStore }

func NewRetryDispatchCoordinator(b *store.RetryDispatchStore) *RetryDispatchCoordinator {
	return &RetryDispatchCoordinator{backend: b}
}
func (c *RetryDispatchCoordinator) Send(key string) error {
	err := c.backend.Attempt(key)
	if err == nil {
		return nil
	}
	wrapped := fmt.Errorf("dispatch failed: %v", err)
	var temporary *store.RetryDispatchTemporaryError
	if errors.As(wrapped, &temporary) {
		return c.backend.Attempt(key)
	}
	return wrapped
}
