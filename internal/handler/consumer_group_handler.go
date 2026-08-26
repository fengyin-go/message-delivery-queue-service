package handler

import (
	"net/http"

	"message-queue/internal/model"
	"message-queue/pkg/httpx"
)

func (s *Server) registerConsumerGroupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/consumer-groups", s.createConsumerGroup)
	mux.HandleFunc("GET /api/consumer-groups", s.listConsumerGroups)
	mux.HandleFunc("GET /api/consumer-groups/{id}", s.getConsumerGroup)
	mux.HandleFunc("PUT /api/consumer-groups/{id}", s.updateConsumerGroup)
	mux.HandleFunc("DELETE /api/consumer-groups/{id}", s.deleteConsumerGroup)
	mux.HandleFunc("POST /api/consumer-groups/{id}/advance", s.advanceOffset)
}

type createConsumerGroupRequest struct {
	Name    string `json:"name"`
	TopicID string `json:"topic_id"`
}

func (s *Server) createConsumerGroup(w http.ResponseWriter, r *http.Request) {
	var req createConsumerGroupRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	cg, err := s.svc.CreateConsumerGroup(model.ConsumerGroup{Name: req.Name, TopicID: req.TopicID})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, cg)
}

func (s *Server) listConsumerGroups(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ConsumerGroupFilter{
		TopicID: r.URL.Query().Get("topic_id"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListConsumerGroups(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getConsumerGroup(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	cg, err := s.svc.GetConsumerGroup(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, cg)
}

func (s *Server) updateConsumerGroup(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createConsumerGroupRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	cg, err := s.svc.UpdateConsumerGroup(id, model.ConsumerGroup{Name: req.Name, TopicID: req.TopicID})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, cg)
}

func (s *Server) deleteConsumerGroup(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteConsumerGroup(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type advanceOffsetRequest struct {
	Delta int64 `json:"delta"`
}

func (s *Server) advanceOffset(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req advanceOffsetRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	cg, err := s.svc.AdvanceOffset(id, req.Delta)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, cg)
}
