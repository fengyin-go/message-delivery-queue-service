package store

import "message-queue/internal/model"

func (s *MemoryStore) CreateConsumerGroup(cg *model.ConsumerGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.consumerGroups {
		if exist.Name == cg.Name {
			return ErrConflict
		}
	}
	s.consumerGroups[cg.ID] = cg
	return nil
}

func (s *MemoryStore) GetConsumerGroup(id string) (*model.ConsumerGroup, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cg, ok := s.consumerGroups[id]
	if !ok {
		return nil, ErrNotFound
	}
	return cg, nil
}

func (s *MemoryStore) GetConsumerGroupByName(name string) (*model.ConsumerGroup, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, cg := range s.consumerGroups {
		if cg.Name == name {
			return cg, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListConsumerGroups() []*model.ConsumerGroup {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.ConsumerGroup, 0, len(s.consumerGroups))
	for _, cg := range s.consumerGroups {
		list = append(list, cg)
	}
	return list
}

func (s *MemoryStore) UpdateConsumerGroup(cg *model.ConsumerGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.consumerGroups[cg.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.consumerGroups {
		if exist.ID != cg.ID && exist.Name == cg.Name {
			return ErrConflict
		}
	}
	s.consumerGroups[cg.ID] = cg
	return nil
}

func (s *MemoryStore) DeleteConsumerGroup(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.consumerGroups[id]; !ok {
		return ErrNotFound
	}
	delete(s.consumerGroups, id)
	return nil
}
