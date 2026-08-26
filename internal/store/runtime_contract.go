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
// Build assembles the config for key. The candidate is published to the shared
// cache only once assembly — including the failure-prone payload decode — has
// completed successfully. A failed build panics before publishing anything, so
// no later Get can ever observe a half-initialized Assembly.
func (s *TopicConfigAssemblyStore) Build(key string, fail bool) *Assembly {
	candidate := &Assembly{ID: key}
	if fail {
		panic("payload decode")
	}
	candidate.Parts = map[string]string{"header": "ready"}
	s.mu.Lock()
	s.items[key] = candidate
	s.mu.Unlock()
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
