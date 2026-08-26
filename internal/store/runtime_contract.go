package store

import (
	"errors"
	"sync"
)

type ConsumerLeaseBatchPool struct {
	mu          sync.Mutex
	open, limit int
}
type ConsumerLeaseBatchLease struct {
	pool *ConsumerLeaseBatchPool
	once sync.Once
}

func NewConsumerLeaseBatchPool(limit int) *ConsumerLeaseBatchPool {
	return &ConsumerLeaseBatchPool{limit: limit}
}
func (p *ConsumerLeaseBatchPool) Acquire() (*ConsumerLeaseBatchLease, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.open++
	if p.open > p.limit {
		return nil, errors.New("lease limit")
	}
	return &ConsumerLeaseBatchLease{pool: p}, nil
}
func (l *ConsumerLeaseBatchLease) Close() {
	l.once.Do(func() { l.pool.mu.Lock(); l.pool.open--; l.pool.mu.Unlock() })
}
func (p *ConsumerLeaseBatchPool) Open() int { p.mu.Lock(); defer p.mu.Unlock(); return p.open }
