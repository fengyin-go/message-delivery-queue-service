package service_test

import (
	"message-queue/internal/service"
	"message-queue/internal/store"
	"testing"
)

func TestR025Target(t *testing.T) {
	pool := store.NewAuditWriterLeasePool(2)
	coordinator := service.NewAuditWriterLeaseCoordinator(pool)
	processed, err := coordinator.Process([]string{"a", "b", "c", "d"})
	if err != nil || processed != 4 || pool.Open() != 0 {
		t.Fatalf("批处理在中途耗尽租约并且没有释放已占用资源")
	}
}
