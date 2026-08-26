// Package handler 实现 HTTP 处理器层。
package handler

import (
	"errors"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"message-queue/internal/config"
	"message-queue/internal/model"
	"message-queue/internal/service"
	"message-queue/internal/store"
	"message-queue/pkg/httpx"
	"message-queue/pkg/logger"
)

type Server struct {
	svc *service.Service
	log *logger.Logger
	cfg *config.Config
}

func NewServer(svc *service.Service, log *logger.Logger, cfg *config.Config) *Server {
	return &Server{svc: svc, log: log, cfg: cfg}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	s.registerTopicRoutes(mux)
	s.registerMessageRoutes(mux)
	s.registerProducerRoutes(mux)
	s.registerConsumerRoutes(mux)
	s.registerConsumerGroupRoutes(mux)
	s.registerSubscriptionRoutes(mux)
	s.registerRetryRoutes(mux)
	s.registerDeadLetterRoutes(mux)
	s.registerStatsRoutes(mux)
	s.registerHealthRoutes(mux)
	s.registerBatchRoutes(mux)
	s.registerAuditRoutes(mux)
	return s.rateLimitMiddleware(s.authMiddleware(s.loggingMiddleware(s.recoveryMiddleware(mux))))
}

func (s *Server) maxPageSize() int {
	if s.cfg != nil && s.cfg.MaxPageSize > 0 {
		return s.cfg.MaxPageSize
	}
	return 100
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		s.log.Infof("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func (s *Server) recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				s.log.Errorf("panic: %v\n%s", rec, debug.Stack())
				httpx.InternalError(w, "服务器内部错误")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

const bearerToken = "mq-demo-token"

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || strings.HasPrefix(r.URL.Path, "/frontend/") || strings.HasSuffix(r.URL.Path, ".html") {
			next.ServeHTTP(w, r)
			return
		}
		auth := r.Header.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") || strings.TrimPrefix(auth, "Bearer ") != bearerToken {
			httpx.Unauthorized(w, "未授权")
			return
		}
		next.ServeHTTP(w, r)
	})
}

type tokenBucket struct {
	tokens float64
	last   time.Time
	mu     sync.Mutex
	rate   float64
	cap    float64
}

var defaultBucket = &tokenBucket{tokens: 100, last: time.Now(), rate: 100, cap: 200}

func (b *tokenBucket) allow() bool {
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

func (s *Server) rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !defaultBucket.allow() {
			httpx.Error(w, http.StatusTooManyRequests, 429, "请求过于频繁")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case model.IsValidationError(err):
		httpx.BadRequest(w, err.Error())
	case errors.Is(err, store.ErrNotFound):
		httpx.NotFound(w, err.Error())
	case errors.Is(err, store.ErrConflict):
		httpx.Conflict(w, err.Error())
	default:
		httpx.InternalError(w, err.Error())
	}
}
