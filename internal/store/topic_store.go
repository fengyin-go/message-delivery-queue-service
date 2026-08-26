package store

import "message-queue/internal/model"

func (s *MemoryStore) CreateTopic(t *model.Topic) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.topics {
		if exist.Name == t.Name {
			return ErrConflict
		}
	}
	s.topics[t.ID] = t
	return nil
}

func (s *MemoryStore) GetTopic(id string) (*model.Topic, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.topics[id]
	if !ok {
		return nil, ErrNotFound
	}
	return t, nil
}

func (s *MemoryStore) GetTopicByName(name string) (*model.Topic, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, t := range s.topics {
		if t.Name == name {
			return t, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListTopics() []*model.Topic {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Topic, 0, len(s.topics))
	for _, t := range s.topics {
		list = append(list, t)
	}
	return list
}

func (s *MemoryStore) UpdateTopic(t *model.Topic) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.topics[t.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.topics {
		if exist.ID != t.ID && exist.Name == t.Name {
			return ErrConflict
		}
	}
	s.topics[t.ID] = t
	return nil
}

func (s *MemoryStore) DeleteTopic(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.topics[id]; !ok {
		return ErrNotFound
	}
	delete(s.topics, id)
	return nil
}
