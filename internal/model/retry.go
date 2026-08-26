package model

import (
	"strings"
	"time"
)

const MaxRetryAttempts = 3

type Retry struct {
	ID           string    `json:"id"`
	MessageID    string    `json:"message_id"`
	Attempt      int       `json:"attempt"`
	Reason       string    `json:"reason"`
	NextRetryAt  time.Time `json:"next_retry_at"`
	CreatedAt    time.Time `json:"created_at"`
}

func (r *Retry) Validate() error {
	if r.MessageID == "" {
		return NewValidationError("message_id", "消息 ID 不能为空")
	}
	if r.Attempt <= 0 {
		return NewValidationError("attempt", "重试次数必须大于 0")
	}
	r.Reason = strings.TrimSpace(r.Reason)
	if r.Reason == "" {
		return NewValidationError("reason", "重试原因不能为空")
	}
	if r.NextRetryAt.IsZero() {
		return NewValidationError("next_retry_at", "下次重试时间不能为空")
	}
	return nil
}

type RetryFilter struct {
	MessageID string
}

func (f RetryFilter) Match(r *Retry) bool {
	if f.MessageID != "" && r.MessageID != f.MessageID {
		return false
	}
	return true
}
