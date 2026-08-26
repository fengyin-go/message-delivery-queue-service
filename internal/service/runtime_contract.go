package service

import (
	"context"
	"message-queue/internal/store"
)

type SubscriptionSyncCancelCoordinator struct {
	backend *store.SubscriptionSyncCancelStore
}

func NewSubscriptionSyncCancelCoordinator(b *store.SubscriptionSyncCancelStore) *SubscriptionSyncCancelCoordinator {
	return &SubscriptionSyncCancelCoordinator{backend: b}
}
func (c *SubscriptionSyncCancelCoordinator) Dispatch(ctx context.Context) error {
	for attempt := 0; attempt < 3; attempt++ {
		if err := c.backend.Attempt(context.Background()); err != nil {
			return err
		}
	}
	return nil
}
