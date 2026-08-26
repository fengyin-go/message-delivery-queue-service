package service

import "message-queue/internal/store"

type PayloadArchiveCoordinator struct {
	backend *store.PayloadArchiveStore
	pending []byte
}

func NewPayloadArchiveCoordinator(b *store.PayloadArchiveStore) *PayloadArchiveCoordinator {
	return &PayloadArchiveCoordinator{backend: b}
}
func (c *PayloadArchiveCoordinator) Archive(payload []byte) {
	c.backend.Put(payload)
	c.pending = payload
}
func (c *PayloadArchiveCoordinator) Export() []byte { return append([]byte(nil), c.pending...) }
