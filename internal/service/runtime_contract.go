package service

import "message-queue/internal/store"

type PartitionResultStreamCoordinator struct {
	backend *store.PartitionResultStreamStore
}

func NewPartitionResultStreamCoordinator(b *store.PartitionResultStreamStore) *PartitionResultStreamCoordinator {
	return &PartitionResultStreamCoordinator{backend: b}
}
func (c *PartitionResultStreamCoordinator) Collect(fail bool) (values []string, err error) {
	results, errs := c.backend.Stream(fail)
	for value := range results {
		values = append(values, value)
	}
	if err := <-errs; err != nil {
		return values, err
	}
	return values, nil
}
