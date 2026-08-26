package middleware

import (
	"context"
	"net/http"

	"message-queue/pkg/idgen"
)

// requestIDKey 请求 ID 上下文键类型。
type requestIDKey struct{}

// RequestIDHeader 请求 ID Header 名称。
const RequestIDHeader = "X-Request-ID"

// RequestIDMiddleware 为每个请求附加唯一请求 ID。
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get(RequestIDHeader)
		if rid == "" {
			rid = idgen.Hex()
		}
		w.Header().Set(RequestIDHeader, rid)
		ctx := context.WithValue(r.Context(), requestIDKey{}, rid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRequestID 从上下文中获取请求 ID。
func GetRequestID(ctx context.Context) string {
	if rid, ok := ctx.Value(requestIDKey{}).(string); ok {
		return rid
	}
	return ""
}

// RequestIDMiddlewareWithGenerator 使用自定义生成器的请求 ID 中间件。
func RequestIDMiddlewareWithGenerator(gen func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rid := r.Header.Get(RequestIDHeader)
			if rid == "" {
				rid = gen()
			}
			w.Header().Set(RequestIDHeader, rid)
			ctx := context.WithValue(r.Context(), requestIDKey{}, rid)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ResponseWriterWithRequestID 包装 ResponseWriter 以透传请求 ID。
type ResponseWriterWithRequestID struct {
	http.ResponseWriter
	RequestID string
}

// WriteHeader 写入状态码并附加请求 ID。
func (w *ResponseWriterWithRequestID) WriteHeader(code int) {
	w.ResponseWriter.Header().Set(RequestIDHeader, w.RequestID)
	w.ResponseWriter.WriteHeader(code)
}
