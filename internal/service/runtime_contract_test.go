package service_test

import (
	"message-queue/internal/service"
	"message-queue/internal/store"
	"testing"
)

func TestR024Target(t *testing.T) {
	backend := store.NewProducerSendRetryStore()
	coordinator := service.NewProducerSendRetryCoordinator(backend)
	err := coordinator.Send("message-1")
	records := backend.Records()
	if err != nil || len(records) != 1 || records[0] != "message-1" {
		t.Fatalf("临时失败没有正确重试或留下了重复投递记录")
	}
}
