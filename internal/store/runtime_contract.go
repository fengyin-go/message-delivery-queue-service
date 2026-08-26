package store

import "sync"

type VersionedState struct {
	Version int
	State   string
}
type AcknowledgementVersionStore struct {
	mu     sync.Mutex
	states map[string]VersionedState
}

func NewAcknowledgementVersionStore() *AcknowledgementVersionStore {
	return &AcknowledgementVersionStore{states: make(map[string]VersionedState)}
}
func (s *AcknowledgementVersionStore) Update(key string, version int, state string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.states[key] = VersionedState{Version: version, State: state}
	return true
}
func (s *AcknowledgementVersionStore) Get(key string) VersionedState {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.states[key]
}
