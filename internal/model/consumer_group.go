package model

import (
	"strings"
	"time"
)

type ConsumerGroup struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	TopicID   string    `json:"topic_id"`
	Offset    int64     `json:"offset"`
	CreatedAt time.Time `json:"created_at"`
}

func (cg *ConsumerGroup) Validate() error {
	cg.Name = strings.TrimSpace(cg.Name)
	if cg.Name == "" {
		return NewValidationError("name", "消费组名称不能为空")
	}
	if cg.TopicID == "" {
		return NewValidationError("topic_id", "所属主题不能为空")
	}
	if cg.Offset < 0 {
		return NewValidationError("offset", "offset 不能为负数")
	}
	return nil
}

type ConsumerGroupFilter struct {
	TopicID string
	Keyword string
}

func (f ConsumerGroupFilter) Match(cg *ConsumerGroup) bool {
	if f.TopicID != "" && cg.TopicID != f.TopicID {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(cg.Name), k) {
			return false
		}
	}
	return true
}
