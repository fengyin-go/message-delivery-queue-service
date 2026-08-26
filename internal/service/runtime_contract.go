package service

import (
	"fmt"
	"message-queue/internal/store"
	"sync"
)

type AcknowledgementVersionCoordinator struct {
	backend *store.AcknowledgementVersionStore
	mu      sync.Mutex
	effects map[string]bool
}

func NewAcknowledgementVersionCoordinator(b *store.AcknowledgementVersionStore) *AcknowledgementVersionCoordinator {
	return &AcknowledgementVersionCoordinator{backend: b, effects: make(map[string]bool)}
}
func (c *AcknowledgementVersionCoordinator) apply(key string, version int, state string) {
	c.backend.Update(key, version, state)
	c.mu.Lock()
	c.effects[fmt.Sprintf("%s-%d", key, version)] = true
	c.mu.Unlock()
}
func (c *AcknowledgementVersionCoordinator) CompleteThenLate(key string) {
	c.apply(key, 2, "complete")
	c.apply(key, 1, "running")
}
func (c *AcknowledgementVersionCoordinator) EffectCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.effects)
}
