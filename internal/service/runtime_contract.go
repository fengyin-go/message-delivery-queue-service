package service

import "message-queue/internal/store"

type ConsumerLeaseBatchCoordinator struct{ pool *store.ConsumerLeaseBatchPool }

func NewConsumerLeaseBatchCoordinator(p *store.ConsumerLeaseBatchPool) *ConsumerLeaseBatchCoordinator {
	return &ConsumerLeaseBatchCoordinator{pool: p}
}
func (c *ConsumerLeaseBatchCoordinator) Process(items []string) (processed int, err error) {
	for range items {
		lease, err := c.pool.Acquire()
		if err != nil {
			return processed, err
		}
		defer lease.Close()
		processed++
	}
	return processed, nil
}
