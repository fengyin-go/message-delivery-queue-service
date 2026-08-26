// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"message-queue/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	// Topic
	CreateTopic(t *model.Topic) error
	GetTopic(id string) (*model.Topic, error)
	GetTopicByName(name string) (*model.Topic, error)
	ListTopics() []*model.Topic
	UpdateTopic(t *model.Topic) error
	DeleteTopic(id string) error

	// Message
	CreateMessage(m *model.Message) error
	GetMessage(id string) (*model.Message, error)
	ListMessages() []*model.Message
	UpdateMessage(m *model.Message) error
	DeleteMessage(id string) error
	GetMessagesByTopic(topicID string) []*model.Message

	// Producer
	CreateProducer(p *model.Producer) error
	GetProducer(id string) (*model.Producer, error)
	GetProducerByName(name string) (*model.Producer, error)
	ListProducers() []*model.Producer
	UpdateProducer(p *model.Producer) error
	DeleteProducer(id string) error

	// Consumer
	CreateConsumer(c *model.Consumer) error
	GetConsumer(id string) (*model.Consumer, error)
	GetConsumerByName(name string) (*model.Consumer, error)
	ListConsumers() []*model.Consumer
	UpdateConsumer(c *model.Consumer) error
	DeleteConsumer(id string) error

	// ConsumerGroup
	CreateConsumerGroup(cg *model.ConsumerGroup) error
	GetConsumerGroup(id string) (*model.ConsumerGroup, error)
	GetConsumerGroupByName(name string) (*model.ConsumerGroup, error)
	ListConsumerGroups() []*model.ConsumerGroup
	UpdateConsumerGroup(cg *model.ConsumerGroup) error
	DeleteConsumerGroup(id string) error

	// Subscription
	CreateSubscription(s *model.Subscription) error
	GetSubscription(id string) (*model.Subscription, error)
	ListSubscriptions() []*model.Subscription
	UpdateSubscription(s *model.Subscription) error
	DeleteSubscription(id string) error
	FindSubscription(groupID, topicID string) (*model.Subscription, error)

	// Retry
	CreateRetry(r *model.Retry) error
	GetRetry(id string) (*model.Retry, error)
	ListRetries() []*model.Retry
	ListRetriesByMessage(messageID string) []*model.Retry
	UpdateRetry(r *model.Retry) error
	DeleteRetry(id string) error

	// DeadLetter
	CreateDeadLetter(d *model.DeadLetter) error
	GetDeadLetter(id string) (*model.DeadLetter, error)
	ListDeadLetters() []*model.DeadLetter
	UpdateDeadLetter(d *model.DeadLetter) error
	DeleteDeadLetter(id string) error
	GetDeadLettersByTopic(topicID string) []*model.DeadLetter

	// Audit
	CreateAudit(a *model.Audit) error
	GetAudit(id string) (*model.Audit, error)
	ListAudits() []*model.Audit
	DeleteAudit(id string) error
	GetAuditsByEntity(entityType, entityID string) []*model.Audit
}
