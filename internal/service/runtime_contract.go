package service

import (
	"errors"
	"fmt"
	"message-queue/internal/store"
)

type DeadLetterReplayCoordinator struct{ backend *store.DeadLetterReplayStore }

func NewDeadLetterReplayCoordinator(b *store.DeadLetterReplayStore) *DeadLetterReplayCoordinator {
	return &DeadLetterReplayCoordinator{backend: b}
}
func (c *DeadLetterReplayCoordinator) Send(key string) error {
	err := c.backend.Attempt(key)
	if err == nil {
		return nil
	}
	wrapped := fmt.Errorf("dispatch failed: %v", err)
	var temporary *store.DeadLetterReplayTemporaryError
	if errors.As(wrapped, &temporary) {
		return c.backend.Attempt(key)
	}
	return wrapped
}
