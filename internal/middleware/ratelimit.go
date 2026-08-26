package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"message-queue/pkg/httpx"
)

// TokenBucket 令牌桶限流器。
type TokenBucket struct {
	tokens float64
	last   time.Time
	mu     sync.Mutex
	rate   float64
	cap    float64
}

// NewTokenBucket 创建令牌桶，rate 为每秒产生令牌数，cap 为桶容量。
func NewTokenBucket(rate, cap float64) *TokenBucket {
	return &TokenBucket{
		tokens: cap,
		last:   time.Now(),
		rate:   rate,
		cap:    cap,
	}
}

// Allow 判断是否允许当前请求通过。
func (b *TokenBucket) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	now := time.Now()
	elapsed := now.Sub(b.last).Seconds()
	b.last = now
	b.tokens += elapsed * b.rate
	if b.tokens > b.cap {
		b.tokens = b.cap
	}
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

// RateLimitMiddleware 限流中间件。
func RateLimitMiddleware(bucket *TokenBucket) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !bucket.Allow() {
				httpx.Error(w, http.StatusTooManyRequests, 429, "请求过于频繁")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RateLimitByIPMiddleware 基于 IP 的限流中间件。
type RateLimitByIPMiddleware struct {
	buckets map[string]*TokenBucket
	mu      sync.RWMutex
	rate    float64
	cap     float64
}

// NewRateLimitByIPMiddleware 创建基于 IP 的限流器。
func NewRateLimitByIPMiddleware(rate, cap float64) *RateLimitByIPMiddleware {
	return &RateLimitByIPMiddleware{
		buckets: make(map[string]*TokenBucket),
		rate:    rate,
		cap:     cap,
	}
}

func (m *RateLimitByIPMiddleware) getBucket(ip string) *TokenBucket {
	m.mu.RLock()
	b, ok := m.buckets[ip]
	m.mu.RUnlock()
	if ok {
		return b
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	b, ok = m.buckets[ip]
	if ok {
		return b
	}
	b = NewTokenBucket(m.rate, m.cap)
	m.buckets[ip] = b
	return b
}

// Handler 返回限流中间件。
func (m *RateLimitByIPMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if idx := strings.Index(ip, ":"); idx != -1 {
			ip = ip[:idx]
		}
		if ip == "" {
			ip = "unknown"
		}
		if !m.getBucket(ip).Allow() {
			httpx.Error(w, http.StatusTooManyRequests, 429, "请求过于频繁")
			return
		}
		next.ServeHTTP(w, r)
	})
}
