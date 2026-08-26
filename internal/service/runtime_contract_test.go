package service_test

import (
	"context"
	"errors"
	"message-queue/internal/service"
	"message-queue/internal/store"
	"testing"
	"time"
)

func TestR020Target(t *testing.T) {
	backend := store.NewDeadLetterReplayCancelStore(20 * time.Millisecond)
	coordinator := service.NewDeadLetterReplayCancelCoordinator(backend)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()
	started := time.Now()
	err := coordinator.Dispatch(ctx)
	if !errors.Is(err, context.DeadlineExceeded) || time.Since(started) > 35*time.Millisecond || backend.Calls() != 1 {
		t.Fatalf("请求取消后后台仍继续发起后续投递")
	}
}
