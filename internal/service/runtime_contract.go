package service

import (
	"context"
	"message-queue/internal/store"
)

type ProducerCancellationCoordinator struct {
	backend *store.ProducerCancellationStore
}

func NewProducerCancellationCoordinator(b *store.ProducerCancellationStore) *ProducerCancellationCoordinator {
	return &ProducerCancellationCoordinator{backend: b}
}
func (c *ProducerCancellationCoordinator) Dispatch(ctx context.Context) error {
	for attempt := 0; attempt < 3; attempt++ {
		if err := c.backend.Attempt(context.Background()); err != nil {
			return err
		}
	}
	return nil
}
