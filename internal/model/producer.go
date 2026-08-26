package model

import (
	"strings"
	"time"
)

type Producer struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	TopicID   string    `json:"topic_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (p *Producer) Validate() error {
	p.Name = strings.TrimSpace(p.Name)
	if p.Name == "" {
		return NewValidationError("name", "生产者名称不能为空")
	}
	if p.TopicID == "" {
		return NewValidationError("topic_id", "所属主题不能为空")
	}
	return nil
}

type ProducerFilter struct {
	TopicID string
	Keyword string
}

func (f ProducerFilter) Match(p *Producer) bool {
	if f.TopicID != "" && p.TopicID != f.TopicID {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(p.Name), k) {
			return false
		}
	}
	return true
}
