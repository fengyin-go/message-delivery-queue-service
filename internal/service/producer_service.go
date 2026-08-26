package service

import (
	"sort"
	"time"

	"message-queue/internal/model"
	"message-queue/pkg/idgen"
)

func (s *Service) CreateProducer(input model.Producer) (*model.Producer, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetTopic(input.TopicID); err != nil {
		return nil, model.NewValidationError("topic_id", "所属主题不存在")
	}
	p := &model.Producer{
		ID:        idgen.Hex(),
		Name:      input.Name,
		TopicID:   input.TopicID,
		CreatedAt: time.Now(),
	}
	if err := s.store.CreateProducer(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) GetProducer(id string) (*model.Producer, error) {
	return s.store.GetProducer(id)
}

func (s *Service) ListProducers(filter model.ProducerFilter, page, size int) ([]*model.Producer, int, error) {
	all := s.store.ListProducers()
	matched := make([]*model.Producer, 0, len(all))
	for _, p := range all {
		if filter.Match(p) {
			matched = append(matched, p)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Producer{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateProducer(id string, input model.Producer) (*model.Producer, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	p, err := s.store.GetProducer(id)
	if err != nil {
		return nil, err
	}
	if p.TopicID != input.TopicID {
		if _, err := s.store.GetTopic(input.TopicID); err != nil {
			return nil, model.NewValidationError("topic_id", "所属主题不存在")
		}
	}
	p.Name = input.Name
	p.TopicID = input.TopicID
	if err := s.store.UpdateProducer(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) DeleteProducer(id string) error {
	return s.store.DeleteProducer(id)
}
