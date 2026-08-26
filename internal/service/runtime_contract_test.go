package service_test

import (
	"message-queue/internal/service"
	"message-queue/internal/store"
	"testing"
)

func TestR023Target(t *testing.T) {
	backend := store.NewConsumerTagArchiveStore()
	coordinator := service.NewConsumerTagArchiveCoordinator(backend)
	payload := []byte("alpha")
	coordinator.Archive(payload)
	copy(payload, []byte("omega"))
	if string(backend.Snapshot()) != "alpha" || string(coordinator.Export()) != "alpha" {
		t.Fatalf("首批消息内容被后一次缓冲区复用改写")
	}
}
