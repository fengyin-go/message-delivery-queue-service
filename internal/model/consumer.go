package model

import (
	"strings"
	"time"
)

type Consumer struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func (c *Consumer) Validate() error {
	c.Name = strings.TrimSpace(c.Name)
	if c.Name == "" {
		return NewValidationError("name", "消费者名称不能为空")
	}
	return nil
}

type ConsumerFilter struct {
	Keyword string
}

func (f ConsumerFilter) Match(c *Consumer) bool {
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(c.Name), k) {
			return false
		}
	}
	return true
}
