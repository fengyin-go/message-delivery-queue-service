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
func (c *TopicConfigAssemblyCoordinator) Build(key string, fail bool) (item *store.Assembly, err error) {
	defer func() {
		if recover() != nil {
			item, _ = c.backend.Get(key)
			err = errors.New("assembly failed")
		}
	}()
	item = c.backend.Build(key, fail)
	return item, nil
}
