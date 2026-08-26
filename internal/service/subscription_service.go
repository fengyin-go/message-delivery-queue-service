package service

import (
	"sort"
	"time"

	"message-queue/internal/model"
	"message-queue/pkg/idgen"
)

func (s *Service) CreateSubscription(input model.Subscription) (*model.Subscription, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetConsumerGroup(input.GroupID); err != nil {
		return nil, model.NewValidationError("group_id", "消费组不存在")
	}
	if _, err := s.store.GetTopic(input.TopicID); err != nil {
		return nil, model.NewValidationError("topic_id", "主题不存在")
	}
	if _, err := s.store.FindSubscription(input.GroupID, input.TopicID); err == nil {
		return nil, model.NewValidationError("subscription", "该订阅关系已存在")
	}
	sub := &model.Subscription{
		ID:        idgen.Hex(),
		GroupID:   input.GroupID,
		TopicID:   input.TopicID,
		CreatedAt: time.Now(),
	}
	if err := s.store.CreateSubscription(sub); err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *Service) GetSubscription(id string) (*model.Subscription, error) {
	return s.store.GetSubscription(id)
}

func (s *Service) ListSubscriptions(filter model.SubscriptionFilter, page, size int) ([]*model.Subscription, int, error) {
	all := s.store.ListSubscriptions()
	matched := make([]*model.Subscription, 0, len(all))
	for _, sub := range all {
		if filter.Match(sub) {
			matched = append(matched, sub)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Subscription{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateSubscription(id string, input model.Subscription) (*model.Subscription, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	sub, err := s.store.GetSubscription(id)
	if err != nil {
		return nil, err
	}
	if sub.GroupID != input.GroupID {
		if _, err := s.store.GetConsumerGroup(input.GroupID); err != nil {
			return nil, model.NewValidationError("group_id", "消费组不存在")
		}
	}
	if sub.TopicID != input.TopicID {
		if _, err := s.store.GetTopic(input.TopicID); err != nil {
			return nil, model.NewValidationError("topic_id", "主题不存在")
		}
	}
	sub.GroupID = input.GroupID
	sub.TopicID = input.TopicID
	if err := s.store.UpdateSubscription(sub); err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *Service) DeleteSubscription(id string) error {
	return s.store.DeleteSubscription(id)
}
