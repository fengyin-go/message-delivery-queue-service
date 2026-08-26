// Package middleware 提供 HTTP 中间件。
package middleware

import (
	"net/http"
	"strings"

	"message-queue/pkg/httpx"
)

const BearerToken = "mq-demo-token"

// AuthMiddleware 鉴权中间件，校验 Bearer Token。
func AuthMiddleware(next http.Handler, skipPaths ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, p := range skipPaths {
			if r.URL.Path == p {
				next.ServeHTTP(w, r)
				return
			}
		}
		auth := r.Header.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") ||
			strings.TrimPrefix(auth, "Bearer ") != BearerToken {
			httpx.Unauthorized(w, "未授权")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// AuthFunc 返回指定 token 的鉴权中间件。
func AuthFunc(token string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" || !strings.HasPrefix(auth, "Bearer ") ||
				strings.TrimPrefix(auth, "Bearer ") != token {
				httpx.Unauthorized(w, "未授权")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
