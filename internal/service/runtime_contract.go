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
		// 在启动新一轮同步前检查取消信号，避免在上一轮结束与本轮启动之间
		// 的取消无法阻止本轮发起请求。
		if err := ctx.Err(); err != nil {
			return err
		}
		// 将调用方的 ctx 透传到底层 Attempt，使取消信号能中断正在进行的请求，
		// 而非在不可取消的 background context 上空等到 delay 结束。
		if err := c.backend.Attempt(ctx); err != nil {
			return err
		}
	}
	return nil
}
