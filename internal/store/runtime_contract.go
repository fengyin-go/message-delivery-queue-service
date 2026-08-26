package store

import "sync"

type Assembly struct {
	ID    string
	Parts map[string]string
}
type TopicConfigAssemblyStore struct {
	mu    sync.Mutex
	items map[string]*Assembly
}

func NewTopicConfigAssemblyStore() *TopicConfigAssemblyStore {
	return &TopicConfigAssemblyStore{items: make(map[string]*Assembly)}
}
func (s *TopicConfigAssemblyStore) Build(key string, fail bool) *Assembly {
	s.mu.Lock()
	candidate := &Assembly{ID: key}
	s.items[key] = candidate
	s.mu.Unlock()
	if fail {
		panic("payload decode")
	}
	candidate.Parts = map[string]string{"header": "ready"}
	return candidate
}
func (s *TopicConfigAssemblyStore) Get(key string) (*Assembly, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.items[key]
	return item, ok
}
func (s *TopicConfigAssemblyStore) Delete(key string) {
	s.mu.Lock()
	delete(s.items, key)
	s.mu.Unlock()
}
