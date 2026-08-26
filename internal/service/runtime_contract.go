package service

import (
	"context"
	"message-queue/internal/store"
)

type SubscriptionRefreshCoordinator struct {
	backend *store.SubscriptionRefreshStore
}

func NewSubscriptionRefreshCoordinator(b *store.SubscriptionRefreshStore) *SubscriptionRefreshCoordinator {
	return &SubscriptionRefreshCoordinator{backend: b}
}
func (c *SubscriptionRefreshCoordinator) Probe(ctx context.Context, key string) error {
	return c.backend.Wait(context.Background(), key)
}
