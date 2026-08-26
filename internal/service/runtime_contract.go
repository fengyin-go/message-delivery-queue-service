package service

import (
	"fmt"
	"message-queue/internal/store"
	"sync"
)

type SubscriptionTerminalVersionCoordinator struct {
	backend *store.SubscriptionTerminalVersionStore
	mu      sync.Mutex
	effects map[string]bool
}

func NewSubscriptionTerminalVersionCoordinator(b *store.SubscriptionTerminalVersionStore) *SubscriptionTerminalVersionCoordinator {
	return &SubscriptionTerminalVersionCoordinator{backend: b, effects: make(map[string]bool)}
}
func (c *SubscriptionTerminalVersionCoordinator) apply(key string, version int, state string) {
	c.backend.Update(key, version, state)
	c.mu.Lock()
	c.effects[fmt.Sprintf("%s-%d", key, version)] = true
	c.mu.Unlock()
}
func (c *SubscriptionTerminalVersionCoordinator) CompleteThenLate(key string) {
	c.apply(key, 2, "complete")
	c.apply(key, 1, "running")
}
func (c *SubscriptionTerminalVersionCoordinator) EffectCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.effects)
}
