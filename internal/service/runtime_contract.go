package service

import "message-queue/internal/store"

type AuditWriterLeaseCoordinator struct{ pool *store.AuditWriterLeasePool }

func NewAuditWriterLeaseCoordinator(p *store.AuditWriterLeasePool) *AuditWriterLeaseCoordinator {
	return &AuditWriterLeaseCoordinator{pool: p}
}
func (c *AuditWriterLeaseCoordinator) Process(items []string) (processed int, err error) {
	for range items {
		if err := c.processOne(); err != nil {
			return processed, err
		}
		processed++
	}
	return processed, nil
}

// processOne acquires and releases a single lease. Defining the work as its own
// function scopes the lease's defer to a single iteration, so each lease is
// released as soon as that item is done — not held until the whole batch
// returns. Without this, a batch larger than the pool limit exhausts the pool
// and leaves open leases behind on early exit.
func (c *AuditWriterLeaseCoordinator) processOne() error {
	lease, err := c.pool.Acquire()
	if err != nil {
		return err
	}
	defer lease.Close()
	return nil
}
