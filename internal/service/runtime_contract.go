package service

import "message-queue/internal/store"

type ProducerConnectionLeaseCoordinator struct {
	pool *store.ProducerConnectionLeasePool
}

func NewProducerConnectionLeaseCoordinator(p *store.ProducerConnectionLeasePool) *ProducerConnectionLeaseCoordinator {
	return &ProducerConnectionLeaseCoordinator{pool: p}
}
func (c *ProducerConnectionLeaseCoordinator) Process(items []string) (processed int, err error) {
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
