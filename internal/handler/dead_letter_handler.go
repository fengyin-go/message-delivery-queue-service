package handler

import (
	"net/http"

	"message-queue/internal/model"
	"message-queue/pkg/httpx"
)

func (s *Server) registerDeadLetterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/dead-letters", s.createDeadLetter)
	mux.HandleFunc("GET /api/dead-letters", s.listDeadLetters)
	mux.HandleFunc("GET /api/dead-letters/{id}", s.getDeadLetter)
	mux.HandleFunc("PUT /api/dead-letters/{id}", s.updateDeadLetter)
	mux.HandleFunc("DELETE /api/dead-letters/{id}", s.deleteDeadLetter)
	mux.HandleFunc("GET /api/topics/{topic_id}/dead-letters", s.getDeadLettersByTopic)
}

type createDeadLetterRequest struct {
	MessageID string `json:"message_id"`
	TopicID   string `json:"topic_id"`
	Reason    string `json:"reason"`
}

func (s *Server) createDeadLetter(w http.ResponseWriter, r *http.Request) {
	var req createDeadLetterRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	d, err := s.svc.CreateDeadLetter(model.DeadLetter{MessageID: req.MessageID, TopicID: req.TopicID, Reason: req.Reason})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, d)
}

func (s *Server) listDeadLetters(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.DeadLetterFilter{
		TopicID: r.URL.Query().Get("topic_id"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListDeadLetters(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getDeadLetter(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	d, err := s.svc.GetDeadLetter(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, d)
}

func (s *Server) updateDeadLetter(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createDeadLetterRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	d, err := s.svc.UpdateDeadLetter(id, model.DeadLetter{MessageID: req.MessageID, TopicID: req.TopicID, Reason: req.Reason})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, d)
}

func (s *Server) deleteDeadLetter(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteDeadLetter(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) getDeadLettersByTopic(w http.ResponseWriter, r *http.Request) {
	topicID := r.PathValue("topic_id")
	items, err := s.svc.GetDeadLettersByTopic(topicID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, items)
}
