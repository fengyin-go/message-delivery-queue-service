package model

import (
	"strings"
	"time"
)

const (
	TopicStatusActive   = "active"
	TopicStatusArchived = "archived"
)

type Topic struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Partitions      int       `json:"partitions"`
	RetentionSeconds int      `json:"retention_seconds"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

func (t *Topic) Validate() error {
	t.Name = strings.TrimSpace(t.Name)
	if t.Name == "" {
		return NewValidationError("name", "主题名称不能为空")
	}
	if t.Partitions <= 0 {
		return NewValidationError("partitions", "分区数必须大于 0")
	}
	if t.RetentionSeconds <= 0 {
		return NewValidationError("retention_seconds", "保留时间必须大于 0")
	}
	if t.Status == "" {
		t.Status = TopicStatusActive
	}
	if t.Status != TopicStatusActive && t.Status != TopicStatusArchived {
		return NewValidationError("status", "主题状态不合法")
	}
	return nil
}

type TopicFilter struct {
	Status  string
	Keyword string
}

func (f TopicFilter) Match(t *Topic) bool {
	if f.Status != "" && t.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(t.Name), k) {
			return false
		}
	}
	return true
}
