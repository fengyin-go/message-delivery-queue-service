package store

import "message-queue/internal/model"

func (s *MemoryStore) CreateProducer(p *model.Producer) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.producers {
		if exist.Name == p.Name {
			return ErrConflict
		}
	}
	s.producers[p.ID] = p
	return nil
}

func (s *MemoryStore) GetProducer(id string) (*model.Producer, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.producers[id]
	if !ok {
		return nil, ErrNotFound
	}
	return p, nil
}

func (s *MemoryStore) GetProducerByName(name string) (*model.Producer, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, p := range s.producers {
		if p.Name == name {
			return p, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListProducers() []*model.Producer {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Producer, 0, len(s.producers))
	for _, p := range s.producers {
		list = append(list, p)
	}
	return list
}

func (s *MemoryStore) UpdateProducer(p *model.Producer) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.producers[p.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.producers {
		if exist.ID != p.ID && exist.Name == p.Name {
			return ErrConflict
		}
	}
	s.producers[p.ID] = p
	return nil
}

func (s *MemoryStore) DeleteProducer(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.producers[id]; !ok {
		return ErrNotFound
	}
	delete(s.producers, id)
	return nil
}
