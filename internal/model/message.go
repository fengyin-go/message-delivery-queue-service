package model

import (
	"strings"
	"time"
)

const (
	MessageStatusPending       = "pending"
	MessageStatusDelivered     = "delivered"
	MessageStatusAcknowledged  = "acknowledged"
	MessageStatusDead          = "dead"
)

var messageTransitions = map[string]map[string]bool{
	MessageStatusPending:      {MessageStatusDelivered: true, MessageStatusDead: true},
	MessageStatusDelivered:    {MessageStatusAcknowledged: true, MessageStatusDead: true},
	MessageStatusAcknowledged: {},
	MessageStatusDead:         {},
}

func CanTransitionMessage(from, to string) bool {
	if m, ok := messageTransitions[from]; ok {
		return m[to]
	}
	return false
}

type Message struct {
	ID        string    `json:"id"`
	TopicID   string    `json:"topic_id"`
	Partition int       `json:"partition"`
	Key       string    `json:"key"`
	Payload   string    `json:"payload"`
	Status    string    `json:"status"`
	Attempts  int       `json:"attempts"`
	CreatedAt time.Time `json:"created_at"`
}

func (m *Message) Validate() error {
	m.Key = strings.TrimSpace(m.Key)
	if m.TopicID == "" {
		return NewValidationError("topic_id", "所属主题不能为空")
	}
	if m.Payload == "" {
		return NewValidationError("payload", "消息内容不能为空")
	}
	if m.Partition < 0 {
		return NewValidationError("partition", "分区号不能为负数")
	}
	if m.Status == "" {
		m.Status = MessageStatusPending
	}
	if m.Status != MessageStatusPending && m.Status != MessageStatusDelivered &&
		m.Status != MessageStatusAcknowledged && m.Status != MessageStatusDead {
		return NewValidationError("status", "消息状态不合法")
	}
	return nil
}

type MessageFilter struct {
	TopicID   string
	Status    string
	Partition int
	Keyword   string
}

func (f MessageFilter) Match(m *Message) bool {
	if f.TopicID != "" && m.TopicID != f.TopicID {
		return false
	}
	if f.Status != "" && m.Status != f.Status {
		return false
	}
	if f.Partition >= 0 && m.Partition != f.Partition {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(m.Payload), k) && !strings.Contains(strings.ToLower(m.Key), k) {
			return false
		}
	}
	return true
}
