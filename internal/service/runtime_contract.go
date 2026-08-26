package service

import "message-queue/internal/store"

type RetryReasonExportCoordinator struct {
	backend *store.RetryReasonExportStore
	pending []byte
}

func NewRetryReasonExportCoordinator(b *store.RetryReasonExportStore) *RetryReasonExportCoordinator {
	return &RetryReasonExportCoordinator{backend: b}
}
func (c *RetryReasonExportCoordinator) Archive(payload []byte) {
	c.backend.Put(payload)
	c.pending = payload
}
func (c *RetryReasonExportCoordinator) Export() []byte { return append([]byte(nil), c.pending...) }
