package service

import "message-queue/internal/store"

type ConsumerTagArchiveCoordinator struct {
	backend *store.ConsumerTagArchiveStore
	pending []byte
}

func NewConsumerTagArchiveCoordinator(b *store.ConsumerTagArchiveStore) *ConsumerTagArchiveCoordinator {
	return &ConsumerTagArchiveCoordinator{backend: b}
}
// Archive 写入消费者标签归档。归档内容在写入时即与调用方切片脱钩，
// 存档完成后调用方复用原切片不会改写缓存与导出中的已归档标签。
func (c *ConsumerTagArchiveCoordinator) Archive(payload []byte) {
	c.backend.Put(payload)
	c.pending = append([]byte(nil), payload...)
}

// Export 返回归档内容的副本，调用方修改返回值不影响归档状态。
func (c *ConsumerTagArchiveCoordinator) Export() []byte {
	return append([]byte(nil), c.pending...)
}
