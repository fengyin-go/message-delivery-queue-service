package service

import (
	"message-queue/internal/store"
	"runtime"
)

type GroupOffsetSnapshotCoordinator struct {
	backend *store.GroupOffsetSnapshotStore
}

func NewGroupOffsetSnapshotCoordinator(b *store.GroupOffsetSnapshotStore) *GroupOffsetSnapshotCoordinator {
	return &GroupOffsetSnapshotCoordinator{backend: b}
}
func (c *GroupOffsetSnapshotCoordinator) CountDuringUpdate() int {
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
