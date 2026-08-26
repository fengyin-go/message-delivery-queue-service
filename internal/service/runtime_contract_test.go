package service_test

import (
	"message-queue/internal/service"
	"message-queue/internal/store"
	"testing"
)

func TestR028Target(t *testing.T) {
	backend := store.NewTopicConfigAssemblyStore()
	coordinator := service.NewTopicConfigAssemblyCoordinator(backend)
	item, err := coordinator.Build("entry-1", true)
	cached, exists := backend.Get("entry-1")
	if err == nil || item != nil || exists || cached != nil {
		t.Fatalf("首次组装失败后仍缓存了半成品并交给后续读取")
	}
}
