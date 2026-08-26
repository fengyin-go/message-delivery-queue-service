package service

import (
	"sort"
	"time"

	"message-queue/internal/model"
	"message-queue/pkg/idgen"
)

func (s *Service) CreateTopic(input model.Topic) (*model.Topic, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	t := &model.Topic{
		ID:               idgen.Hex(),
		Name:             input.Name,
		Partitions:       input.Partitions,
		RetentionSeconds: input.RetentionSeconds,
		Status:           input.Status,
		CreatedAt:        time.Now(),
	}
	if t.Status == "" {
		t.Status = model.TopicStatusActive
	}
	if err := s.store.CreateTopic(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) GetTopic(id string) (*model.Topic, error) {
	return s.store.GetTopic(id)
}

func (s *Service) ListTopics(filter model.TopicFilter, page, size int) ([]*model.Topic, int, error) {
	all := s.store.ListTopics()
	matched := make([]*model.Topic, 0, len(all))
	for _, t := range all {
		if filter.Match(t) {
			matched = append(matched, t)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Topic{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateTopic(id string, input model.Topic) (*model.Topic, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	t, err := s.store.GetTopic(id)
	if err != nil {
		return nil, err
	}
	t.Name = input.Name
	t.Partitions = input.Partitions
	t.RetentionSeconds = input.RetentionSeconds
	t.Status = input.Status
	if err := s.store.UpdateTopic(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) DeleteTopic(id string) error {
	return s.store.DeleteTopic(id)
}
