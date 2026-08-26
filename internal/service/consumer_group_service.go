package service

import (
	"sort"
	"time"

	"message-queue/internal/model"
	"message-queue/pkg/idgen"
)

func (s *Service) CreateConsumerGroup(input model.ConsumerGroup) (*model.ConsumerGroup, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetTopic(input.TopicID); err != nil {
		return nil, model.NewValidationError("topic_id", "所属主题不存在")
	}
	cg := &model.ConsumerGroup{
		ID:        idgen.Hex(),
		Name:      input.Name,
		TopicID:   input.TopicID,
		Offset:    0,
		CreatedAt: time.Now(),
	}
	if err := s.store.CreateConsumerGroup(cg); err != nil {
		return nil, err
	}
	return cg, nil
}

func (s *Service) GetConsumerGroup(id string) (*model.ConsumerGroup, error) {
	return s.store.GetConsumerGroup(id)
}

func (s *Service) ListConsumerGroups(filter model.ConsumerGroupFilter, page, size int) ([]*model.ConsumerGroup, int, error) {
	all := s.store.ListConsumerGroups()
	matched := make([]*model.ConsumerGroup, 0, len(all))
	for _, cg := range all {
		if filter.Match(cg) {
			matched = append(matched, cg)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.ConsumerGroup{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateConsumerGroup(id string, input model.ConsumerGroup) (*model.ConsumerGroup, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	cg, err := s.store.GetConsumerGroup(id)
	if err != nil {
		return nil, err
	}
	if cg.TopicID != input.TopicID {
		if _, err := s.store.GetTopic(input.TopicID); err != nil {
			return nil, model.NewValidationError("topic_id", "所属主题不存在")
		}
	}
	cg.Name = input.Name
	cg.TopicID = input.TopicID
	cg.Offset = input.Offset
	if err := s.store.UpdateConsumerGroup(cg); err != nil {
		return nil, err
	}
	return cg, nil
}

func (s *Service) DeleteConsumerGroup(id string) error {
	return s.store.DeleteConsumerGroup(id)
}

func (s *Service) AdvanceOffset(id string, delta int64) (*model.ConsumerGroup, error) {
	cg, err := s.store.GetConsumerGroup(id)
	if err != nil {
		return nil, err
	}
	cg.Offset += delta
	if cg.Offset < 0 {
		cg.Offset = 0
	}
	if err := s.store.UpdateConsumerGroup(cg); err != nil {
		return nil, err
	}
	return cg, nil
}
