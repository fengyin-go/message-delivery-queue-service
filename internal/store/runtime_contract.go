package store

import (
	"errors"
	"sync"
)

type ProducerConnectionLeasePool struct {
	mu          sync.Mutex
	open, limit int
}
type ProducerConnectionLeaseLease struct {
	pool *ProducerConnectionLeasePool
	once sync.Once
}

func NewProducerConnectionLeasePool(limit int) *ProducerConnectionLeasePool {
	return &ProducerConnectionLeasePool{limit: limit}
}
func (p *ProducerConnectionLeasePool) Acquire() (*ProducerConnectionLeaseLease, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.open++
	if p.open > p.limit {
		return nil, errors.New("lease limit")
	}
	return &ProducerConnectionLeaseLease{pool: p}, nil
}
func (l *ProducerConnectionLeaseLease) Close() {
	l.once.Do(func() { l.pool.mu.Lock(); l.pool.open--; l.pool.mu.Unlock() })
}
func (p *ProducerConnectionLeasePool) Open() int { p.mu.Lock(); defer p.mu.Unlock(); return p.open }
