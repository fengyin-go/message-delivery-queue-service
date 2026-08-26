package handler

import (
	"net/http"
	"time"

	"message-queue/internal/model"
	"message-queue/internal/scheduler"
	"message-queue/pkg/httpx"
)

func (s *Server) registerHealthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", s.healthCheck)
	mux.HandleFunc("GET /health/ready", s.readinessCheck)
	mux.HandleFunc("GET /health/live", s.livenessCheck)
}

type healthResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Checks    map[string]string `json:"checks,omitempty"`
}

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	resp := healthResponse{
		Status:    "ok",
		Timestamp: time.Now().Format(time.RFC3339),
		Checks: map[string]string{
			"store": "up",
		},
	}
	httpx.OK(w, resp)
}

func (s *Server) readinessCheck(w http.ResponseWriter, r *http.Request) {
	_, _, err := s.svc.ListTopics(model.TopicFilter{}, 1, 1)
	if err != nil {
		httpx.Error(w, http.StatusServiceUnavailable, 503, "not ready")
		return
	}
	httpx.OK(w, healthResponse{Status: "ready", Timestamp: time.Now().Format(time.RFC3339)})
}

func (s *Server) livenessCheck(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, healthResponse{Status: "alive", Timestamp: time.Now().Format(time.RFC3339)})
}

// RegisterHealthSchedulerRoutes 注册健康调度器相关路由。
func RegisterHealthSchedulerRoutes(mux *http.ServeMux, hs *scheduler.HealthScheduler) {
	mux.HandleFunc("GET /health/detail", func(w http.ResponseWriter, r *http.Request) {
		status := hs.GetStatus()
		httpx.OK(w, status)
	})
}
