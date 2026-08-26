package service

import (
	"sort"
	"time"

	"message-queue/internal/model"
	"message-queue/pkg/idgen"
)

func (s *Service) CreateConsumer(input model.Consumer) (*model.Consumer, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	c := &model.Consumer{
		ID:        idgen.Hex(),
		Name:      input.Name,
		CreatedAt: time.Now(),
	}
	if err := s.store.CreateConsumer(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) GetConsumer(id string) (*model.Consumer, error) {
	return s.store.GetConsumer(id)
}

func (s *Service) ListConsumers(filter model.ConsumerFilter, page, size int) ([]*model.Consumer, int, error) {
	all := s.store.ListConsumers()
	matched := make([]*model.Consumer, 0, len(all))
	for _, c := range all {
		if filter.Match(c) {
			matched = append(matched, c)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Consumer{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateConsumer(id string, input model.Consumer) (*model.Consumer, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	c, err := s.store.GetConsumer(id)
	if err != nil {
		return nil, err
	}
	c.Name = input.Name
	if err := s.store.UpdateConsumer(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) DeleteConsumer(id string) error {
	return s.store.DeleteConsumer(id)
}
