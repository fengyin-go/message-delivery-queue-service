// Package scheduler 提供后台调度任务。
package scheduler

import (
	"time"

	"message-queue/internal/model"
	"message-queue/internal/store"
	"message-queue/pkg/logger"
)

// CleanupScheduler 过期数据清理调度器。
type CleanupScheduler struct {
	store  store.Store
	log    *logger.Logger
	stopCh chan struct{}
}

// NewCleanupScheduler 创建清理调度器。
func NewCleanupScheduler(st store.Store, log *logger.Logger) *CleanupScheduler {
	return &CleanupScheduler{
		store:  st,
		log:    log,
		stopCh: make(chan struct{}),
	}
}

// Start 启动定时清理任务。
func (s *CleanupScheduler) Start(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				s.runCleanup()
			case <-s.stopCh:
				return
			}
		}
	}()
}

// Stop 停止调度器。
func (s *CleanupScheduler) Stop() {
	close(s.stopCh)
}

func (s *CleanupScheduler) runCleanup() {
	s.cleanupExpiredMessages()
	s.cleanupOldRetries()
}

func (s *CleanupScheduler) cleanupExpiredMessages() {
	topics := s.store.ListTopics()
	now := time.Now()
	for _, topic := range topics {
		if topic.Status != model.TopicStatusActive {
			continue
		}
		ttl := time.Duration(topic.RetentionSeconds) * time.Second
		messages := s.store.GetMessagesByTopic(topic.ID)
		for _, msg := range messages {
			if msg.Status == model.MessageStatusAcknowledged && now.Sub(msg.CreatedAt) > ttl {
				if err := s.store.DeleteMessage(msg.ID); err == nil {
					s.log.Debugf("清理过期消息: %s", msg.ID)
				}
			}
		}
	}
}

func (s *CleanupScheduler) cleanupOldRetries() {
	retries := s.store.ListRetries()
	now := time.Now()
	for _, r := range retries {
		if now.After(r.NextRetryAt) && now.Sub(r.CreatedAt) > 24*time.Hour {
			if err := s.store.DeleteRetry(r.ID); err == nil {
				s.log.Debugf("清理过期重试记录: %s", r.ID)
			}
		}
	}
}

// CleanupDeadLetters 清理超过保留时间的死信。
func (s *CleanupScheduler) CleanupDeadLetters(maxAge time.Duration) int {
	deadLetters := s.store.ListDeadLetters()
	now := time.Now()
	count := 0
	for _, dl := range deadLetters {
		if now.Sub(dl.MovedAt) > maxAge {
			if err := s.store.DeleteDeadLetter(dl.ID); err == nil {
				count++
			}
		}
	}
	return count
}

// CleanupOrphanProducers 清理引用不存在主题的生产者。
func (s *CleanupScheduler) CleanupOrphanProducers() int {
	producers := s.store.ListProducers()
	count := 0
	for _, p := range producers {
		if _, err := s.store.GetTopic(p.TopicID); err != nil {
			if err := s.store.DeleteProducer(p.ID); err == nil {
				count++
			}
		}
	}
	return count
}

// CleanupOrphanConsumerGroups 清理引用不存在主题的消费组。
func (s *CleanupScheduler) CleanupOrphanConsumerGroups() int {
	groups := s.store.ListConsumerGroups()
	count := 0
	for _, g := range groups {
		if _, err := s.store.GetTopic(g.TopicID); err != nil {
			if err := s.store.DeleteConsumerGroup(g.ID); err == nil {
				count++
			}
		}
	}
	return count
}

// CleanupOrphanSubscriptions 清理引用不存在消费组或主题的订阅。
func (s *CleanupScheduler) CleanupOrphanSubscriptions() int {
	subs := s.store.ListSubscriptions()
	count := 0
	for _, sub := range subs {
		_, err1 := s.store.GetConsumerGroup(sub.GroupID)
		_, err2 := s.store.GetTopic(sub.TopicID)
		if err1 != nil || err2 != nil {
			if err := s.store.DeleteSubscription(sub.ID); err == nil {
				count++
			}
		}
	}
	return count
}
