package service_test

import (
	"message-queue/internal/service"
	"message-queue/internal/store"
	"testing"
	"time"
)

func TestR026Target(t *testing.T) {
	coordinator := service.NewRetryBatchStreamCoordinator(store.NewRetryBatchStreamStore())
	done := make(chan error, 1)
	go func() { _, err := coordinator.Collect(true); done <- err }()
	select {
	case err := <-done:
		if err == nil {
			t.Fatalf("分区失败没有返回错误并结束结果流")
		}
	case <-time.After(80 * time.Millisecond):
		t.Fatalf("分区失败后结果流一直等待且无法结束")
	}
}
