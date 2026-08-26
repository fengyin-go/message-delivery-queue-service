package service

import (
	"sort"
	"time"

	"message-queue/internal/model"
	"message-queue/pkg/idgen"
)

func (s *Service) CreateAudit(input model.Audit) (*model.Audit, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	a := &model.Audit{
		ID:         idgen.Hex(),
		EntityType: input.EntityType,
		EntityID:   input.EntityID,
		Action:     input.Action,
		Operator:   input.Operator,
		Detail:     input.Detail,
		CreatedAt:  time.Now(),
	}
	if err := s.store.CreateAudit(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) GetAudit(id string) (*model.Audit, error) {
	return s.store.GetAudit(id)
}

func (s *Service) ListAudits(filter model.AuditFilter, page, size int) ([]*model.Audit, int, error) {
	all := s.store.ListAudits()
	matched := make([]*model.Audit, 0, len(all))
	for _, a := range all {
		if filter.Match(a) {
			matched = append(matched, a)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Audit{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) DeleteAudit(id string) error {
	return s.store.DeleteAudit(id)
}

func (s *Service) GetAuditsByEntity(entityType, entityID string) ([]*model.Audit, error) {
	return s.store.GetAuditsByEntity(entityType, entityID), nil
}

func (s *Service) RecordAudit(entityType, entityID, action, operator, detail string) (*model.Audit, error) {
	return s.CreateAudit(model.Audit{
		EntityType: entityType,
		EntityID:   entityID,
		Action:     action,
		Operator:   operator,
		Detail:     detail,
	})
}

func (s *Service) RecordMessageAudit(messageID, action, detail string) (*model.Audit, error) {
	return s.RecordAudit("message", messageID, action, "system", detail)
}

func (s *Service) RecordTopicAudit(topicID, action, detail string) (*model.Audit, error) {
	return s.RecordAudit("topic", topicID, action, "system", detail)
}
