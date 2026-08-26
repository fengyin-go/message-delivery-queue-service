package service

import (
	"errors"
	"message-queue/internal/store"
)

type DeadLetterAssemblyCoordinator struct {
	backend *store.DeadLetterAssemblyStore
}

func NewDeadLetterAssemblyCoordinator(b *store.DeadLetterAssemblyStore) *DeadLetterAssemblyCoordinator {
	return &DeadLetterAssemblyCoordinator{backend: b}
}
func (c *DeadLetterAssemblyCoordinator) Build(key string, fail bool) (item *store.Assembly, err error) {
	defer func() {
		if recover() != nil {
			item, _ = c.backend.Get(key)
			err = errors.New("assembly failed")
		}
	}()
	item = c.backend.Build(key, fail)
	return item, nil
}
