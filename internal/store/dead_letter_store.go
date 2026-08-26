package store

import "message-queue/internal/model"

func (s *MemoryStore) CreateDeadLetter(d *model.DeadLetter) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.deadLetters[d.ID] = d
	return nil
}

func (s *MemoryStore) GetDeadLetter(id string) (*model.DeadLetter, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	d, ok := s.deadLetters[id]
	if !ok {
		return nil, ErrNotFound
	}
	return d, nil
}

func (s *MemoryStore) ListDeadLetters() []*model.DeadLetter {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.DeadLetter, 0, len(s.deadLetters))
	for _, d := range s.deadLetters {
		list = append(list, d)
	}
	return list
}

func (s *MemoryStore) UpdateDeadLetter(d *model.DeadLetter) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.deadLetters[d.ID]; !ok {
		return ErrNotFound
	}
	s.deadLetters[d.ID] = d
	return nil
}

func (s *MemoryStore) DeleteDeadLetter(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.deadLetters[id]; !ok {
		return ErrNotFound
	}
	delete(s.deadLetters, id)
	return nil
}

func (s *MemoryStore) GetDeadLettersByTopic(topicID string) []*model.DeadLetter {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.DeadLetter, 0)
	for _, d := range s.deadLetters {
		if d.TopicID == topicID {
			list = append(list, d)
		}
	}
	return list
}
