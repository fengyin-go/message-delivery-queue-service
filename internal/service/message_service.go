package service

import (
	"sort"
	"time"

	"message-queue/internal/model"
	"message-queue/pkg/idgen"
)

func (s *Service) CreateMessage(input model.Message) (*model.Message, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetTopic(input.TopicID); err != nil {
		return nil, model.NewValidationError("topic_id", "所属主题不存在")
	}
	m := &model.Message{
		ID:        idgen.Hex(),
		TopicID:   input.TopicID,
		Partition: input.Partition,
		Key:       input.Key,
		Payload:   input.Payload,
		Status:    model.MessageStatusPending,
		Attempts:  0,
		CreatedAt: time.Now(),
	}
	if err := s.store.CreateMessage(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *Service) GetMessage(id string) (*model.Message, error) {
	return s.store.GetMessage(id)
}

func (s *Service) ListMessages(filter model.MessageFilter, page, size int) ([]*model.Message, int, error) {
	all := s.store.ListMessages()
	matched := make([]*model.Message, 0, len(all))
	for _, m := range all {
		if filter.Match(m) {
			matched = append(matched, m)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Message{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateMessage(id string, input model.Message) (*model.Message, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	m, err := s.store.GetMessage(id)
	if err != nil {
		return nil, err
	}
	m.TopicID = input.TopicID
	m.Partition = input.Partition
	m.Key = input.Key
	m.Payload = input.Payload
	if err := s.store.UpdateMessage(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *Service) DeleteMessage(id string) error {
	return s.store.DeleteMessage(id)
}

func (s *Service) DeliverMessage(id string) (*model.Message, error) {
	m, err := s.store.GetMessage(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionMessage(m.Status, model.MessageStatusDelivered) {
		return nil, model.NewValidationError("status", "非法状态流转")
	}
	m.Status = model.MessageStatusDelivered
	m.Attempts++
	if err := s.store.UpdateMessage(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *Service) AckMessage(id string, groupID string) (*model.Message, error) {
	m, err := s.store.GetMessage(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionMessage(m.Status, model.MessageStatusAcknowledged) {
		return nil, model.NewValidationError("status", "非法状态流转")
	}
	m.Status = model.MessageStatusAcknowledged
	if err := s.store.UpdateMessage(m); err != nil {
		return nil, err
	}
	if groupID != "" {
		cg, err := s.store.GetConsumerGroup(groupID)
		if err == nil && cg != nil {
			cg.Offset++
			_ = s.store.UpdateConsumerGroup(cg)
		}
	}
	return m, nil
}

func (s *Service) FailMessage(id string, reason string) (*model.Message, error) {
	m, err := s.store.GetMessage(id)
	if err != nil {
		return nil, err
	}
	if m.Status != model.MessageStatusPending && m.Status != model.MessageStatusDelivered {
		return nil, model.NewValidationError("status", "当前状态不允许标记失败")
	}
	m.Attempts++
	if m.Attempts >= model.MaxRetryAttempts {
		m.Status = model.MessageStatusDead
		if err := s.store.UpdateMessage(m); err != nil {
			return nil, err
		}
		dl := &model.DeadLetter{
			ID:        idgen.Hex(),
			MessageID: m.ID,
			TopicID:   m.TopicID,
			Reason:    reason,
			MovedAt:   time.Now(),
			CreatedAt: time.Now(),
		}
		_ = s.store.CreateDeadLetter(dl)
		return m, nil
	}
	retry := &model.Retry{
		ID:          idgen.Hex(),
		MessageID:   m.ID,
		Attempt:     m.Attempts,
		Reason:      reason,
		NextRetryAt: time.Now().Add(time.Duration(m.Attempts) * time.Minute),
		CreatedAt:   time.Now(),
	}
	_ = s.store.CreateRetry(retry)
	if err := s.store.UpdateMessage(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *Service) BatchCreateMessages(inputs []model.Message) ([]*model.Message, error) {
	results := make([]*model.Message, 0, len(inputs))
	for _, input := range inputs {
		m, err := s.CreateMessage(input)
		if err != nil {
			return nil, err
		}
		results = append(results, m)
	}
	return results, nil
}

func (s *Service) GetMessagesByTopic(topicID string) ([]*model.Message, error) {
	if _, err := s.store.GetTopic(topicID); err != nil {
		return nil, model.NewValidationError("topic_id", "所属主题不存在")
	}
	return s.store.GetMessagesByTopic(topicID), nil
}
