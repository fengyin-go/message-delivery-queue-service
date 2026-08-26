package scheduler

import (
	"sync"
	"time"

	"message-queue/internal/store"
	"message-queue/pkg/logger"
)

// HealthStatus 健康状态。
type HealthStatus struct {
	Healthy   bool      `json:"healthy"`
	CheckedAt time.Time `json:"checked_at"`
	Message   string    `json:"message,omitempty"`
}

// HealthScheduler 健康检查调度器。
type HealthScheduler struct {
	store    store.Store
	log      *logger.Logger
	stopCh   chan struct{}
	status   HealthStatus
	statusMu sync.RWMutex
	interval time.Duration
}

// NewHealthScheduler 创建健康检查调度器。
func NewHealthScheduler(st store.Store, log *logger.Logger) *HealthScheduler {
	return &HealthScheduler{
		store:    st,
		log:      log,
		stopCh:   make(chan struct{}),
		interval: 30 * time.Second,
		status:   HealthStatus{Healthy: true, CheckedAt: time.Now()},
	}
}

// Start 启动健康检查。
func (s *HealthScheduler) Start() {
	ticker := time.NewTicker(s.interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				s.check()
			case <-s.stopCh:
				return
			}
		}
	}()
}

// Stop 停止健康检查。
func (s *HealthScheduler) Stop() {
	close(s.stopCh)
}

// GetStatus 获取当前健康状态。
func (s *HealthScheduler) GetStatus() HealthStatus {
	s.statusMu.RLock()
	defer s.statusMu.RUnlock()
	return s.status
}

func (s *HealthScheduler) check() {
	now := time.Now()
	s.statusMu.Lock()
	defer s.statusMu.Unlock()

	topics := s.store.ListTopics()
	if topics == nil {
		s.status = HealthStatus{Healthy: false, CheckedAt: now, Message: "store unavailable"}
		s.log.Warnf("健康检查失败: topics nil")
		return
	}
	messages := s.store.ListMessages()
	if messages == nil {
		s.status = HealthStatus{Healthy: false, CheckedAt: now, Message: "store unavailable"}
		s.log.Warnf("健康检查失败: messages nil")
		return
	}

	s.status = HealthStatus{Healthy: true, CheckedAt: now, Message: "ok"}
}

// CheckNow 立即执行一次健康检查。
func (s *HealthScheduler) CheckNow() HealthStatus {
	s.check()
	return s.GetStatus()
}

// SetInterval 设置检查间隔。
func (s *HealthScheduler) SetInterval(d time.Duration) {
	s.interval = d
}
