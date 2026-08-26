package service_test

import (
	"message-queue/internal/service"
	"message-queue/internal/store"
	"testing"
)

func TestR029Target(t *testing.T) {
	coordinator := service.NewGroupOffsetSnapshotCoordinator(store.NewGroupOffsetSnapshotStore(1))
	if coordinator.CountDuringUpdate() != 1 {
		t.Fatalf("采集中的主题快照被同时到达的更新改写")
	}
}
