package service

import (
	"fmt"
	"message-queue/internal/store"
	"sync"
)

type MessageDeliveryVersionCoordinator struct {
	backend *store.MessageDeliveryVersionStore
	mu      sync.Mutex
	effects map[string]bool
}

func NewMessageDeliveryVersionCoordinator(b *store.MessageDeliveryVersionStore) *MessageDeliveryVersionCoordinator {
	return &MessageDeliveryVersionCoordinator{backend: b, effects: make(map[string]bool)}
}
// apply 按版本号提交一次投递终态，并仅在写入被接受时记录一次副作用。
// store 层会拒绝版本倒退/同版本重放，因此迟到回调不会覆盖终态，
// 也不会因重放产生重复副作用（如重复外部确认）。
func (c *MessageDeliveryVersionCoordinator) apply(key string, version int, state string) {
	if !c.backend.Update(key, version, state) {
		// 旧回调晚到或重放：终态已由更高版本提交，丢弃，不产生副作用。
		return
	}
	c.mu.Lock()
	c.effects[fmt.Sprintf("%s-%d", key, version)] = true
	c.mu.Unlock()
}
// CompleteThenLate 复现“终态先提交、首轮延迟回调后到”的竞态：
// 先以版本 2 提交 complete 终态，再重放版本 1 的 running 回调。
// 修复后迟到回调（版本 1<=2）被 store 拒绝，终态保持 complete，
// 副作用计数应为 1，而非被覆盖回 running 且确认两次。
func (c *MessageDeliveryVersionCoordinator) CompleteThenLate(key string) {
	c.apply(key, 2, "complete")
	c.apply(key, 1, "running")
}
func (c *MessageDeliveryVersionCoordinator) EffectCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.effects)
}
