package store

import (
	"errors"
	"sync"
)

type AuditWriterLeasePool struct {
	mu          sync.Mutex
	open, limit int
}
type AuditWriterLeaseLease struct {
	pool *AuditWriterLeasePool
	once sync.Once
}

func NewAuditWriterLeasePool(limit int) *AuditWriterLeasePool {
	return &AuditWriterLeasePool{limit: limit}
}
func (p *AuditWriterLeasePool) Acquire() (*AuditWriterLeaseLease, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.open >= p.limit {
		return nil, errors.New("lease limit")
	}
	p.open++
	return &AuditWriterLeaseLease{pool: p}, nil
}
func (l *AuditWriterLeaseLease) Close() {
	l.once.Do(func() { l.pool.mu.Lock(); l.pool.open--; l.pool.mu.Unlock() })
}
func (p *AuditWriterLeasePool) Open() int { p.mu.Lock(); defer p.mu.Unlock(); return p.open }
