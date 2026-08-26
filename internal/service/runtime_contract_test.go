package service_test

import (
	"message-queue/internal/service"
	"message-queue/internal/store"
	"testing"
)

func TestR027Target(t *testing.T) {
	backend := store.NewMessageDeliveryVersionStore()
	coordinator := service.NewMessageDeliveryVersionCoordinator(backend)
	coordinator.CompleteThenLate("delivery-1")
	state := backend.Get("delivery-1")
	if state.State != "complete" || state.Version != 2 || coordinator.EffectCount() != 1 {
		t.Fatalf("延迟回调覆盖了新终态或重复执行了确认副作用")
	}
}
