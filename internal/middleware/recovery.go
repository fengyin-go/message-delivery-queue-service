package middleware

import (
	"net/http"
	"runtime/debug"

	"message-queue/pkg/httpx"
	"message-queue/pkg/logger"
)

// RecoveryMiddleware panic 恢复中间件。
func RecoveryMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log.Errorf("panic: %v\n%s", rec, debug.Stack())
					httpx.InternalError(w, "服务器内部错误")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// RecoveryMiddlewareWithHook 带自定义回调的 panic 恢复中间件。
func RecoveryMiddlewareWithHook(
	log *logger.Logger,
	hook func(rec interface{}, stack []byte),
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					stack := debug.Stack()
					log.Errorf("panic: %v\n%s", rec, stack)
					if hook != nil {
						hook(rec, stack)
					}
					httpx.InternalError(w, "服务器内部错误")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
