package middleware

import "net/http"

// Chain 将多个中间件串联成一个。
func Chain(middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(final http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}

// ChainHandlers 将多个 http.Handler 按顺序执行（用于降级场景）。
func ChainHandlers(handlers ...http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, h := range handlers {
			h.ServeHTTP(w, r)
		}
	})
}

// WrapHandler 将 http.HandlerFunc 包装为中间件形式。
func WrapHandler(fn http.HandlerFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fn(w, r)
			next.ServeHTTP(w, r)
		})
	}
}

// If 条件中间件，仅在条件满足时应用。
func If(condition func(r *http.Request) bool, mw func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if condition(r) {
				mw(next).ServeHTTP(w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// Skip 跳过指定路径的中间件。
func Skip(paths []string, mw func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, p := range paths {
				if r.URL.Path == p {
					next.ServeHTTP(w, r)
					return
				}
			}
			mw(next).ServeHTTP(w, r)
		})
	}
}

// Compose 组合中间件（从左到右应用）。
func Compose(mws ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		result := h
		for _, mw := range mws {
			result = mw(result)
		}
		return result
	}
}
