package service

import (
	"context"
	"message-queue/internal/store"
)

type DeadLetterReplayCancelCoordinator struct {
	backend *store.DeadLetterReplayCancelStore
}

func NewDeadLetterReplayCancelCoordinator(b *store.DeadLetterReplayCancelStore) *DeadLetterReplayCancelCoordinator {
	return &DeadLetterReplayCancelCoordinator{backend: b}
}
func (c *DeadLetterReplayCancelCoordinator) Dispatch(ctx context.Context) error {
	for attempt := 0; attempt < 3; attempt++ {
		if err := c.backend.Attempt(context.Background()); err != nil {
			return err
		}
	}
	return nil
}
