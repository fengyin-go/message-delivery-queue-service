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
		// 超时/取消后不再安排下一轮投递。
		if err := ctx.Err(); err != nil {
			return err
		}
		// 将调用方 ctx 透传给下游，使取消信号能及时中断当前投递。
		if err := c.backend.Attempt(ctx); err != nil {
			return err
		}
	}
	return nil
}
