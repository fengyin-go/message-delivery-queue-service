package service

import (
	"context"
	"message-queue/internal/store"
)

type TimedMessageProbeCoordinator struct{ backend *store.TimedMessageProbeStore }

func NewTimedMessageProbeCoordinator(b *store.TimedMessageProbeStore) *TimedMessageProbeCoordinator {
	return &TimedMessageProbeCoordinator{backend: b}
}
func (c *TimedMessageProbeCoordinator) Probe(ctx context.Context, key string) error {
	return c.backend.Wait(context.Background(), key)
}
