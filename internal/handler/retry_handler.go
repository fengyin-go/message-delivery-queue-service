package handler

import (
	"net/http"
	"time"

	"message-queue/internal/model"
	"message-queue/pkg/httpx"
)

func (s *Server) registerRetryRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/retries", s.createRetry)
	mux.HandleFunc("GET /api/retries", s.listRetries)
	mux.HandleFunc("GET /api/retries/{id}", s.getRetry)
	mux.HandleFunc("PUT /api/retries/{id}", s.updateRetry)
	mux.HandleFunc("DELETE /api/retries/{id}", s.deleteRetry)
	mux.HandleFunc("GET /api/messages/{message_id}/retries", s.getRetriesByMessage)
}

type createRetryRequest struct {
	MessageID   string    `json:"message_id"`
	Attempt     int       `json:"attempt"`
	Reason      string    `json:"reason"`
	NextRetryAt time.Time `json:"next_retry_at"`
}

func (s *Server) createRetry(w http.ResponseWriter, r *http.Request) {
	var req createRetryRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	retry, err := s.svc.CreateRetry(model.Retry{MessageID: req.MessageID, Attempt: req.Attempt, Reason: req.Reason, NextRetryAt: req.NextRetryAt})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, retry)
}

func (s *Server) listRetries(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.RetryFilter{
		MessageID: r.URL.Query().Get("message_id"),
	}
	items, total, err := s.svc.ListRetries(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getRetry(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	retry, err := s.svc.GetRetry(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, retry)
}

func (s *Server) updateRetry(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createRetryRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	retry, err := s.svc.UpdateRetry(id, model.Retry{MessageID: req.MessageID, Attempt: req.Attempt, Reason: req.Reason, NextRetryAt: req.NextRetryAt})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, retry)
}

func (s *Server) deleteRetry(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteRetry(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) getRetriesByMessage(w http.ResponseWriter, r *http.Request) {
	messageID := r.PathValue("message_id")
	items, err := s.svc.GetRetriesByMessage(messageID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, items)
}
