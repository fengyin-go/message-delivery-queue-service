package service

import (
	"errors"
	"message-queue/internal/store"
)

type TopicConfigAssemblyCoordinator struct {
	backend *store.TopicConfigAssemblyStore
}

func NewTopicConfigAssemblyCoordinator(b *store.TopicConfigAssemblyStore) *TopicConfigAssemblyCoordinator {
	return &TopicConfigAssemblyCoordinator{backend: b}
}
// Build assembles the config for key, surfacing assembly failures as ordinary
// errors. On failure the backend panics before publishing anything, so there
// is no cache entry to recover — the shared cache must never expose a
// half-initialized Assembly to a subsequent request.
func (c *TopicConfigAssemblyCoordinator) Build(key string, fail bool) (item *store.Assembly, err error) {
	defer func() {
		if r := recover(); r != nil {
			item = nil
			err = errors.New("assembly failed")
		}
	}()
	item = c.backend.Build(key, fail)
	return item, nil
}
