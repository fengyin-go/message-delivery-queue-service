package store

import "message-queue/internal/model"

func (s *MemoryStore) CreateMessage(m *model.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messages[m.ID] = m
	return nil
}

func (s *MemoryStore) GetMessage(id string) (*model.Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.messages[id]
	if !ok {
		return nil, ErrNotFound
	}
	return m, nil
}

func (s *MemoryStore) ListMessages() []*model.Message {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Message, 0, len(s.messages))
	for _, m := range s.messages {
		list = append(list, m)
	}
	return list
}

func (s *MemoryStore) UpdateMessage(m *model.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.messages[m.ID]; !ok {
		return ErrNotFound
	}
	s.messages[m.ID] = m
	return nil
}

func (s *MemoryStore) DeleteMessage(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.messages[id]; !ok {
		return ErrNotFound
	}
	delete(s.messages, id)
	return nil
}

func (s *MemoryStore) GetMessagesByTopic(topicID string) []*model.Message {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Message, 0)
	for _, m := range s.messages {
		if m.TopicID == topicID {
			list = append(list, m)
		}
	}
	return list
}
