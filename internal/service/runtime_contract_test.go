package service_test

import (
	"context"
	"errors"
	"message-queue/internal/service"
	"message-queue/internal/store"
	"testing"
	"time"
)

func TestR021Target(t *testing.T) {
	backend := store.NewConsumerSessionProbeStore(45 * time.Millisecond)
	coordinator := service.NewConsumerSessionProbeCoordinator(backend)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()
	started := time.Now()
	firstErr := coordinator.Probe(ctx, "first")
	elapsed := time.Since(started)
	secondErr := coordinator.Probe(context.Background(), "second")
	if !errors.Is(firstErr, context.DeadlineExceeded) || elapsed > 30*time.Millisecond || secondErr != nil {
		t.Fatalf("取消的请求没有及时停止或污染了后续正常请求")
	}
}
