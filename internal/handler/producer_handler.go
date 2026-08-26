package handler

import (
	"net/http"

	"message-queue/internal/model"
	"message-queue/pkg/httpx"
)

func (s *Server) registerProducerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/producers", s.createProducer)
	mux.HandleFunc("GET /api/producers", s.listProducers)
	mux.HandleFunc("GET /api/producers/{id}", s.getProducer)
	mux.HandleFunc("PUT /api/producers/{id}", s.updateProducer)
	mux.HandleFunc("DELETE /api/producers/{id}", s.deleteProducer)
}

type createProducerRequest struct {
	Name    string `json:"name"`
	TopicID string `json:"topic_id"`
}

func (s *Server) createProducer(w http.ResponseWriter, r *http.Request) {
	var req createProducerRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	p, err := s.svc.CreateProducer(model.Producer{Name: req.Name, TopicID: req.TopicID})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, p)
}

func (s *Server) listProducers(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ProducerFilter{
		TopicID: r.URL.Query().Get("topic_id"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListProducers(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getProducer(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	p, err := s.svc.GetProducer(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}

func (s *Server) updateProducer(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createProducerRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	p, err := s.svc.UpdateProducer(id, model.Producer{Name: req.Name, TopicID: req.TopicID})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}

func (s *Server) deleteProducer(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteProducer(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
