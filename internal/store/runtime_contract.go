package store

import (
	"context"
	"time"
)

// ConsumerSessionProbeStore 维护消费者会话探测的等待窗口。
//
// 每次探测只受自身调用的生命周期控制：调用方的 ctx 取消则立即返回，
// 否则等待 delay 时长后正常结束。不同探测之间不共享会话状态，因此
// 前一次会话超时或取消不会波及后续探测。
type ConsumerSessionProbeStore struct {
	delay time.Duration
}

func NewConsumerSessionProbeStore(delay time.Duration) *ConsumerSessionProbeStore {
	return &ConsumerSessionProbeStore{delay: delay}
}

// Wait 阻塞至本次探测的等待窗口结束，或在调用方 ctx 被取消时立即返回。
// key 标识探测目标，供上层路由/日志使用；等待本身按调用独立，不跨会话共享。
func (s *ConsumerSessionProbeStore) Wait(ctx context.Context, key string) error {
	timer := time.NewTimer(s.delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
