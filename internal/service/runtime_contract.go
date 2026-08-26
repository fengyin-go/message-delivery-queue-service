package service

import (
	"errors"
	"message-queue/internal/store"
)

type AuditEntryAssemblyCoordinator struct {
	backend *store.AuditEntryAssemblyStore
}

func NewAuditEntryAssemblyCoordinator(b *store.AuditEntryAssemblyStore) *AuditEntryAssemblyCoordinator {
	return &AuditEntryAssemblyCoordinator{backend: b}
}
func (c *AuditEntryAssemblyCoordinator) Build(key string, fail bool) (item *store.Assembly, err error) {
	defer func() {
		if recover() != nil {
			item, _ = c.backend.Get(key)
			err = errors.New("assembly failed")
		}
	}()
	item = c.backend.Build(key, fail)
	return item, nil
}
