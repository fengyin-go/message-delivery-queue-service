package scheduler

import (
	"sync"
	"time"

	"message-queue/internal/model"
	"message-queue/internal/store"
	"message-queue/pkg/logger"
)

// TopicMetrics 主题指标。
type TopicMetrics struct {
	TopicID         string    `json:"topic_id"`
	PendingCount    int       `json:"pending_count"`
	DeliveredCount  int       `json:"delivered_count"`
	AckedCount      int       `json:"acked_count"`
	DeadCount       int       `json:"dead_count"`
	Throughput      float64   `json:"throughput"`
	RecordedAt      time.Time `json:"recorded_at"`
}

// MonitorScheduler 监控指标调度器。
type MonitorScheduler struct {
	store     store.Store
	log       *logger.Logger
	stopCh    chan struct{}
	metrics   []TopicMetrics
	metricsMu sync.RWMutex
	interval  time.Duration
}

// NewMonitorScheduler 创建监控调度器。
func NewMonitorScheduler(st store.Store, log *logger.Logger) *MonitorScheduler {
	return &MonitorScheduler{
		store:    st,
		log:      log,
		stopCh:   make(chan struct{}),
		interval: 60 * time.Second,
		metrics:  make([]TopicMetrics, 0),
	}
}

// Start 启动监控采集。
func (s *MonitorScheduler) Start() {
	ticker := time.NewTicker(s.interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				s.collect()
			case <-s.stopCh:
				return
			}
		}
	}()
}

// Stop 停止监控采集。
func (s *MonitorScheduler) Stop() {
	close(s.stopCh)
}

// GetMetrics 获取最新指标。
func (s *MonitorScheduler) GetMetrics() []TopicMetrics {
	s.metricsMu.RLock()
	defer s.metricsMu.RUnlock()
	result := make([]TopicMetrics, len(s.metrics))
	copy(result, s.metrics)
	return result
}

func (s *MonitorScheduler) collect() {
	topics := s.store.ListTopics()
	messages := s.store.ListMessages()
	deadLetters := s.store.ListDeadLetters()

	statsMap := make(map[string]*TopicMetrics, len(topics))
	for _, t := range topics {
		statsMap[t.ID] = &TopicMetrics{
			TopicID:    t.ID,
			RecordedAt: time.Now(),
		}
	}

	for _, m := range messages {
		st, ok := statsMap[m.TopicID]
		if !ok {
			continue
		}
		switch m.Status {
		case model.MessageStatusPending:
			st.PendingCount++
		case model.MessageStatusDelivered:
			st.DeliveredCount++
		case model.MessageStatusAcknowledged:
			st.AckedCount++
		case model.MessageStatusDead:
			st.DeadCount++
		}
	}

	for _, dl := range deadLetters {
		st, ok := statsMap[dl.TopicID]
		if ok {
			st.DeadCount++
		}
	}

	for _, st := range statsMap {
		total := st.PendingCount + st.DeliveredCount + st.AckedCount + st.DeadCount
		if total > 0 {
			st.Throughput = float64(st.DeliveredCount+st.AckedCount) / float64(total)
		}
	}

	result := make([]TopicMetrics, 0, len(statsMap))
	for _, st := range statsMap {
		result = append(result, *st)
	}

	s.metricsMu.Lock()
	s.metrics = result
	s.metricsMu.Unlock()

	s.log.Debugf("监控采集完成，共 %d 个主题", len(result))
}

// CollectNow 立即采集一次指标。
func (s *MonitorScheduler) CollectNow() []TopicMetrics {
	s.collect()
	return s.GetMetrics()
}

// SetInterval 设置采集间隔。
func (s *MonitorScheduler) SetInterval(d time.Duration) {
	s.interval = d
}

// GetTopicMetric 获取指定主题的指标。
func (s *MonitorScheduler) GetTopicMetric(topicID string) (TopicMetrics, bool) {
	s.metricsMu.RLock()
	defer s.metricsMu.RUnlock()
	for _, m := range s.metrics {
		if m.TopicID == topicID {
			return m, true
		}
	}
	return TopicMetrics{}, false
}

// GetTotalMessageCount 获取消息总数。
func (s *MonitorScheduler) GetTotalMessageCount() int {
	return len(s.store.ListMessages())
}

// GetTotalDeadLetterCount 获取死信总数。
func (s *MonitorScheduler) GetTotalDeadLetterCount() int {
	return len(s.store.ListDeadLetters())
}

// GetTotalProducerCount 获取生产者总数。
func (s *MonitorScheduler) GetTotalProducerCount() int {
	return len(s.store.ListProducers())
}

// GetTotalConsumerCount 获取消费者总数。
func (s *MonitorScheduler) GetTotalConsumerCount() int {
	return len(s.store.ListConsumers())
}

// GetTotalGroupCount 获取消费组总数。
func (s *MonitorScheduler) GetTotalGroupCount() int {
	return len(s.store.ListConsumerGroups())
}
