package handler

import (
	"net/http"

	"message-queue/internal/model"
	"message-queue/pkg/httpx"
)

func (s *Server) registerSubscriptionRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/subscriptions", s.createSubscription)
	mux.HandleFunc("GET /api/subscriptions", s.listSubscriptions)
	mux.HandleFunc("GET /api/subscriptions/{id}", s.getSubscription)
	mux.HandleFunc("PUT /api/subscriptions/{id}", s.updateSubscription)
	mux.HandleFunc("DELETE /api/subscriptions/{id}", s.deleteSubscription)
}

type createSubscriptionRequest struct {
	GroupID string `json:"group_id"`
	TopicID string `json:"topic_id"`
}

func (s *Server) createSubscription(w http.ResponseWriter, r *http.Request) {
	var req createSubscriptionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sub, err := s.svc.CreateSubscription(model.Subscription{GroupID: req.GroupID, TopicID: req.TopicID})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, sub)
}

func (s *Server) listSubscriptions(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.SubscriptionFilter{
		GroupID: r.URL.Query().Get("group_id"),
		TopicID: r.URL.Query().Get("topic_id"),
	}
	items, total, err := s.svc.ListSubscriptions(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getSubscription(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	sub, err := s.svc.GetSubscription(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sub)
}

func (s *Server) updateSubscription(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createSubscriptionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sub, err := s.svc.UpdateSubscription(id, model.Subscription{GroupID: req.GroupID, TopicID: req.TopicID})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sub)
}

func (s *Server) deleteSubscription(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteSubscription(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
