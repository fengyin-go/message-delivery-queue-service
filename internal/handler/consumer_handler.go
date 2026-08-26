package handler

import (
	"net/http"

	"message-queue/internal/model"
	"message-queue/pkg/httpx"
)

func (s *Server) registerConsumerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/consumers", s.createConsumer)
	mux.HandleFunc("GET /api/consumers", s.listConsumers)
	mux.HandleFunc("GET /api/consumers/{id}", s.getConsumer)
	mux.HandleFunc("PUT /api/consumers/{id}", s.updateConsumer)
	mux.HandleFunc("DELETE /api/consumers/{id}", s.deleteConsumer)
}

type createConsumerRequest struct {
	Name string `json:"name"`
}

func (s *Server) createConsumer(w http.ResponseWriter, r *http.Request) {
	var req createConsumerRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.CreateConsumer(model.Consumer{Name: req.Name})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, c)
}

func (s *Server) listConsumers(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ConsumerFilter{
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListConsumers(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getConsumer(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	c, err := s.svc.GetConsumer(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

func (s *Server) updateConsumer(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createConsumerRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.UpdateConsumer(id, model.Consumer{Name: req.Name})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

func (s *Server) deleteConsumer(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteConsumer(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
