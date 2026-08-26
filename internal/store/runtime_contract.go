package store

import "sync"

type VersionedState struct {
	Version int
	State   string
}
type MessageDeliveryVersionStore struct {
	mu     sync.Mutex
	states map[string]VersionedState
}

func NewMessageDeliveryVersionStore() *MessageDeliveryVersionStore {
	return &MessageDeliveryVersionStore{states: make(map[string]VersionedState)}
}
func (s *MessageDeliveryVersionStore) Update(key string, version int, state string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.states[key] = VersionedState{Version: version, State: state}
	return true
}
func (s *MessageDeliveryVersionStore) Get(key string) VersionedState {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.states[key]
}
