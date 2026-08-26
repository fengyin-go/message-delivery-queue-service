package service

import "message-queue/internal/store"

type ConsumerGroupStreamCoordinator struct {
	backend *store.ConsumerGroupStreamStore
}

func NewConsumerGroupStreamCoordinator(b *store.ConsumerGroupStreamStore) *ConsumerGroupStreamCoordinator {
	return &ConsumerGroupStreamCoordinator{backend: b}
}
func (c *ConsumerGroupStreamCoordinator) Collect(fail bool) (values []string, err error) {
	results, errs := c.backend.Stream(fail)
	for value := range results {
		values = append(values, value)
	}
	if err := <-errs; err != nil {
		return values, err
	}
	return values, nil
}
