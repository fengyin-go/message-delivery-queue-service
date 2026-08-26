package service

import (
	"message-queue/internal/store"
	"runtime"
)

type MessageRateSnapshotCoordinator struct {
	backend *store.MessageRateSnapshotStore
}

func NewMessageRateSnapshotCoordinator(b *store.MessageRateSnapshotStore) *MessageRateSnapshotCoordinator {
	return &MessageRateSnapshotCoordinator{backend: b}
}
func (c *MessageRateSnapshotCoordinator) CountDuringUpdate() int {
	snapshot := c.backend.Snapshot()
	done := make(chan struct{})
	go func() { c.backend.Increment(); close(done) }()
	for i := 0; i < 20000; i++ {
		_ = *snapshot
		runtime.Gosched()
	}
	<-done
	return *snapshot
}
