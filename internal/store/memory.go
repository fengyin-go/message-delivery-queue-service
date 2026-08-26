package store

import (
	"sync"

	"message-queue/internal/model"
)

type MemoryStore struct {
	mu             sync.RWMutex
	topics         map[string]*model.Topic
	messages       map[string]*model.Message
	producers      map[string]*model.Producer
	consumers      map[string]*model.Consumer
	consumerGroups map[string]*model.ConsumerGroup
	subscriptions  map[string]*model.Subscription
	retries        map[string]*model.Retry
	deadLetters    map[string]*model.DeadLetter
	audits         map[string]*model.Audit
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		topics:         make(map[string]*model.Topic),
		messages:       make(map[string]*model.Message),
		producers:      make(map[string]*model.Producer),
		consumers:      make(map[string]*model.Consumer),
		consumerGroups: make(map[string]*model.ConsumerGroup),
		subscriptions:  make(map[string]*model.Subscription),
		retries:        make(map[string]*model.Retry),
		deadLetters:    make(map[string]*model.DeadLetter),
		audits:         make(map[string]*model.Audit),
	}
}

var _ Store = (*MemoryStore)(nil)
