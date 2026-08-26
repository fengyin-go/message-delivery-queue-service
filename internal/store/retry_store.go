package store

import "message-queue/internal/model"

func (s *MemoryStore) CreateRetry(r *model.Retry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.retries[r.ID] = r
	return nil
}

func (s *MemoryStore) GetRetry(id string) (*model.Retry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.retries[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

func (s *MemoryStore) ListRetries() []*model.Retry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Retry, 0, len(s.retries))
	for _, r := range s.retries {
		list = append(list, r)
	}
	return list
}

func (s *MemoryStore) ListRetriesByMessage(messageID string) []*model.Retry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Retry, 0)
	for _, r := range s.retries {
		if r.MessageID == messageID {
			list = append(list, r)
		}
	}
	return list
}

func (s *MemoryStore) UpdateRetry(r *model.Retry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.retries[r.ID]; !ok {
		return ErrNotFound
	}
	s.retries[r.ID] = r
	return nil
}

func (s *MemoryStore) DeleteRetry(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.retries[id]; !ok {
		return ErrNotFound
	}
	delete(s.retries, id)
	return nil
}
