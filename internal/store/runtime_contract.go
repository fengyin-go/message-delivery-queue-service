package store

import "sync"

type VersionedState struct {
	Version int
	State   string
}
type SubscriptionTerminalVersionStore struct {
	mu     sync.Mutex
	states map[string]VersionedState
}

func NewSubscriptionTerminalVersionStore() *SubscriptionTerminalVersionStore {
	return &SubscriptionTerminalVersionStore{states: make(map[string]VersionedState)}
}
func (s *SubscriptionTerminalVersionStore) Update(key string, version int, state string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.states[key] = VersionedState{Version: version, State: state}
	return true
}
func (s *SubscriptionTerminalVersionStore) Get(key string) VersionedState {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.states[key]
}
