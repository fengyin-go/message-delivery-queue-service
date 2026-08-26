package model

import (
	"time"
)

type Subscription struct {
	ID        string    `json:"id"`
	GroupID   string    `json:"group_id"`
	TopicID   string    `json:"topic_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *Subscription) Validate() error {
	if s.GroupID == "" {
		return NewValidationError("group_id", "消费组不能为空")
	}
	if s.TopicID == "" {
		return NewValidationError("topic_id", "主题不能为空")
	}
	return nil
}

type SubscriptionFilter struct {
	GroupID string
	TopicID string
}

func (f SubscriptionFilter) Match(s *Subscription) bool {
	if f.GroupID != "" && s.GroupID != f.GroupID {
		return false
	}
	if f.TopicID != "" && s.TopicID != f.TopicID {
		return false
	}
	return true
}
