package service

import (
	"sort"
	"time"

	"message-queue/internal/model"
	"message-queue/pkg/idgen"
)

func (s *Service) CreateRetry(input model.Retry) (*model.Retry, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetMessage(input.MessageID); err != nil {
		return nil, model.NewValidationError("message_id", "消息不存在")
	}
	r := &model.Retry{
		ID:          idgen.Hex(),
		MessageID:   input.MessageID,
		Attempt:     input.Attempt,
		Reason:      input.Reason,
		NextRetryAt: input.NextRetryAt,
		CreatedAt:   time.Now(),
	}
	if err := s.store.CreateRetry(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) GetRetry(id string) (*model.Retry, error) {
	return s.store.GetRetry(id)
}

func (s *Service) ListRetries(filter model.RetryFilter, page, size int) ([]*model.Retry, int, error) {
	all := s.store.ListRetries()
	matched := make([]*model.Retry, 0, len(all))
	for _, r := range all {
		if filter.Match(r) {
			matched = append(matched, r)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Retry{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateRetry(id string, input model.Retry) (*model.Retry, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	r, err := s.store.GetRetry(id)
	if err != nil {
		return nil, err
	}
	r.MessageID = input.MessageID
	r.Attempt = input.Attempt
	r.Reason = input.Reason
	r.NextRetryAt = input.NextRetryAt
	if err := s.store.UpdateRetry(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) DeleteRetry(id string) error {
	return s.store.DeleteRetry(id)
}

func (s *Service) GetRetriesByMessage(messageID string) ([]*model.Retry, error) {
	return s.store.ListRetriesByMessage(messageID), nil
}
