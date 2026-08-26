package service

import (
	"message-queue/internal/store"
	"runtime"
)

type TopicCounterSnapshotCoordinator struct {
	backend *store.TopicCounterSnapshotStore
}

func NewTopicCounterSnapshotCoordinator(b *store.TopicCounterSnapshotStore) *TopicCounterSnapshotCoordinator {
	return &TopicCounterSnapshotCoordinator{backend: b}
}
func (c *TopicCounterSnapshotCoordinator) CountDuringUpdate() int {
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
