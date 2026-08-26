package middleware

import (
	"net/http"
	"time"

	"message-queue/pkg/logger"
)

// LoggingMiddleware 请求日志中间件。
func LoggingMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			log.Infof("%s %s %s", r.Method, r.URL.Path, time.Since(start))
		})
	}
}

// LoggingMiddlewareWithBody 带请求体摘要的日志中间件。
func LoggingMiddlewareWithBody(log *logger.Logger, maxBodyLen int) func(http.Handler) http.Handler {
	if maxBodyLen <= 0 {
		maxBodyLen = 200
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			log.Infof("%s %s %s", r.Method, r.URL.Path, time.Since(start))
		})
	}
}

// responseRecorder 记录响应状态码。
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

// LoggingMiddlewareWithStatus 带状态码的日志中间件。
func LoggingMiddlewareWithStatus(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rr := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(rr, r)
			log.Infof("%s %s %d %s", r.Method, r.URL.Path, rr.statusCode, time.Since(start))
		})
	}
}
