package store

import "message-queue/internal/model"

func (s *MemoryStore) CreateAudit(a *model.Audit) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.audits[a.ID] = a
	return nil
}

func (s *MemoryStore) GetAudit(id string) (*model.Audit, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.audits[id]
	if !ok {
		return nil, ErrNotFound
	}
	return a, nil
}

func (s *MemoryStore) ListAudits() []*model.Audit {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Audit, 0, len(s.audits))
	for _, a := range s.audits {
		list = append(list, a)
	}
	return list
}

func (s *MemoryStore) DeleteAudit(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.audits[id]; !ok {
		return ErrNotFound
	}
	delete(s.audits, id)
	return nil
}

func (s *MemoryStore) GetAuditsByEntity(entityType, entityID string) []*model.Audit {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Audit, 0)
	for _, a := range s.audits {
		if a.EntityType == entityType && a.EntityID == entityID {
			list = append(list, a)
		}
	}
	return list
}
