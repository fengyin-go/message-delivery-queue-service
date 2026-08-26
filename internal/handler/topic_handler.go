package handler

import (
	"net/http"

	"message-queue/internal/model"
	"message-queue/pkg/httpx"
)

func (s *Server) registerTopicRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/topics", s.createTopic)
	mux.HandleFunc("GET /api/topics", s.listTopics)
	mux.HandleFunc("GET /api/topics/{id}", s.getTopic)
	mux.HandleFunc("PUT /api/topics/{id}", s.updateTopic)
	mux.HandleFunc("DELETE /api/topics/{id}", s.deleteTopic)
}

type createTopicRequest struct {
	Name             string `json:"name"`
	Partitions       int    `json:"partitions"`
	RetentionSeconds int    `json:"retention_seconds"`
	Status           string `json:"status"`
}

func (s *Server) createTopic(w http.ResponseWriter, r *http.Request) {
	var req createTopicRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.CreateTopic(model.Topic{Name: req.Name, Partitions: req.Partitions, RetentionSeconds: req.RetentionSeconds, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, t)
}

func (s *Server) listTopics(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.TopicFilter{
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListTopics(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getTopic(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	t, err := s.svc.GetTopic(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

func (s *Server) updateTopic(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createTopicRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.UpdateTopic(id, model.Topic{Name: req.Name, Partitions: req.Partitions, RetentionSeconds: req.RetentionSeconds, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

func (s *Server) deleteTopic(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteTopic(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
