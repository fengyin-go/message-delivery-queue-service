package service

import (
	"context"
	"message-queue/internal/store"
)

type ConsumerSessionProbeCoordinator struct {
	backend *store.ConsumerSessionProbeStore
}

func NewConsumerSessionProbeCoordinator(b *store.ConsumerSessionProbeStore) *ConsumerSessionProbeCoordinator {
	return &ConsumerSessionProbeCoordinator{backend: b}
}
func (c *ConsumerSessionProbeCoordinator) Probe(ctx context.Context, key string) error {
	return c.backend.Wait(context.Background(), key)
}
