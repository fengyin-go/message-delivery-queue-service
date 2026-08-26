package service

import "message-queue/internal/store"

type AuditWriterLeaseCoordinator struct{ pool *store.AuditWriterLeasePool }

func NewAuditWriterLeaseCoordinator(p *store.AuditWriterLeasePool) *AuditWriterLeaseCoordinator {
	return &AuditWriterLeaseCoordinator{pool: p}
}
func (c *AuditWriterLeaseCoordinator) Process(items []string) (processed int, err error) {
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
