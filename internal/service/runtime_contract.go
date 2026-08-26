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

// Probe 发起一次消费者会话探测。透传调用方 ctx，使每次探测的生命周期
// 由调用方决定（取消立即返回），而非沿用前一个会话的等待状态。
func (c *ConsumerSessionProbeCoordinator) Probe(ctx context.Context, key string) error {
	return c.backend.Wait(ctx, key)
}
