package model

import (
	"strings"
	"time"
)

type DeadLetter struct {
	ID        string    `json:"id"`
	MessageID string    `json:"message_id"`
	TopicID   string    `json:"topic_id"`
	Reason    string    `json:"reason"`
	MovedAt   time.Time `json:"moved_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (d *DeadLetter) Validate() error {
	if d.MessageID == "" {
		return NewValidationError("message_id", "消息 ID 不能为空")
	}
	if d.TopicID == "" {
		return NewValidationError("topic_id", "主题不能为空")
	}
	d.Reason = strings.TrimSpace(d.Reason)
	if d.Reason == "" {
		return NewValidationError("reason", "死信原因不能为空")
	}
	return nil
}

type DeadLetterFilter struct {
	TopicID string
	Keyword string
}

func (f DeadLetterFilter) Match(d *DeadLetter) bool {
	if f.TopicID != "" && d.TopicID != f.TopicID {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(d.Reason), k) {
			return false
		}
	}
	return true
}
