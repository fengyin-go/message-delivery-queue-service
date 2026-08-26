package service

import "message-queue/internal/store"

type ConsumerTagArchiveCoordinator struct {
	backend *store.ConsumerTagArchiveStore
	pending []byte
}

func NewConsumerTagArchiveCoordinator(b *store.ConsumerTagArchiveStore) *ConsumerTagArchiveCoordinator {
	return &ConsumerTagArchiveCoordinator{backend: b}
}
func (c *ConsumerTagArchiveCoordinator) Archive(payload []byte) {
	c.backend.Put(payload)
	c.pending = payload
}
func (c *ConsumerTagArchiveCoordinator) Export() []byte { return append([]byte(nil), c.pending...) }
