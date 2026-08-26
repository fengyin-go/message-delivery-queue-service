package service

import (
	"sort"
	"time"

	"message-queue/internal/model"
	"message-queue/pkg/idgen"
)

func (s *Service) CreateDeadLetter(input model.DeadLetter) (*model.DeadLetter, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetMessage(input.MessageID); err != nil {
		return nil, model.NewValidationError("message_id", "消息不存在")
	}
	if _, err := s.store.GetTopic(input.TopicID); err != nil {
		return nil, model.NewValidationError("topic_id", "主题不存在")
	}
	d := &model.DeadLetter{
		ID:        idgen.Hex(),
		MessageID: input.MessageID,
		TopicID:   input.TopicID,
		Reason:    input.Reason,
		MovedAt:   time.Now(),
		CreatedAt: time.Now(),
	}
	if err := s.store.CreateDeadLetter(d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) GetDeadLetter(id string) (*model.DeadLetter, error) {
	return s.store.GetDeadLetter(id)
}

func (s *Service) ListDeadLetters(filter model.DeadLetterFilter, page, size int) ([]*model.DeadLetter, int, error) {
	all := s.store.ListDeadLetters()
	matched := make([]*model.DeadLetter, 0, len(all))
	for _, d := range all {
		if filter.Match(d) {
			matched = append(matched, d)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.DeadLetter{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateDeadLetter(id string, input model.DeadLetter) (*model.DeadLetter, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	d, err := s.store.GetDeadLetter(id)
	if err != nil {
		return nil, err
	}
	d.MessageID = input.MessageID
	d.TopicID = input.TopicID
	d.Reason = input.Reason
	if err := s.store.UpdateDeadLetter(d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) DeleteDeadLetter(id string) error {
	return s.store.DeleteDeadLetter(id)
}

func (s *Service) GetDeadLettersByTopic(topicID string) ([]*model.DeadLetter, error) {
	if _, err := s.store.GetTopic(topicID); err != nil {
		return nil, model.NewValidationError("topic_id", "主题不存在")
	}
	return s.store.GetDeadLettersByTopic(topicID), nil
}
