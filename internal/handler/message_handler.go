package handler

import (
	"net/http"

	"message-queue/internal/model"
	"message-queue/pkg/httpx"
)

func (s *Server) registerMessageRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/messages", s.createMessage)
	mux.HandleFunc("GET /api/messages", s.listMessages)
	mux.HandleFunc("GET /api/messages/{id}", s.getMessage)
	mux.HandleFunc("PUT /api/messages/{id}", s.updateMessage)
	mux.HandleFunc("DELETE /api/messages/{id}", s.deleteMessage)
	mux.HandleFunc("POST /api/messages/{id}/deliver", s.deliverMessage)
	mux.HandleFunc("POST /api/messages/{id}/ack", s.ackMessage)
	mux.HandleFunc("POST /api/messages/{id}/fail", s.failMessage)
	mux.HandleFunc("POST /api/messages/batch", s.batchCreateMessages)
	mux.HandleFunc("GET /api/topics/{topic_id}/messages", s.getMessagesByTopic)
}

type createMessageRequest struct {
	TopicID   string `json:"topic_id"`
	Partition int    `json:"partition"`
	Key       string `json:"key"`
	Payload   string `json:"payload"`
}

func (s *Server) createMessage(w http.ResponseWriter, r *http.Request) {
	var req createMessageRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	m, err := s.svc.CreateMessage(model.Message{TopicID: req.TopicID, Partition: req.Partition, Key: req.Key, Payload: req.Payload})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, m)
}

func (s *Server) listMessages(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.MessageFilter{
		TopicID:   r.URL.Query().Get("topic_id"),
		Status:    r.URL.Query().Get("status"),
		Keyword:   r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListMessages(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getMessage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	m, err := s.svc.GetMessage(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, m)
}

func (s *Server) updateMessage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createMessageRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	m, err := s.svc.UpdateMessage(id, model.Message{TopicID: req.TopicID, Partition: req.Partition, Key: req.Key, Payload: req.Payload})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, m)
}

func (s *Server) deleteMessage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteMessage(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) deliverMessage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	m, err := s.svc.DeliverMessage(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, m)
}

type ackMessageRequest struct {
	GroupID string `json:"group_id"`
}

func (s *Server) ackMessage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req ackMessageRequest
	_ = httpx.Decode(r, &req)
	m, err := s.svc.AckMessage(id, req.GroupID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, m)
}

type failMessageRequest struct {
	Reason string `json:"reason"`
}

func (s *Server) failMessage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req failMessageRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	m, err := s.svc.FailMessage(id, req.Reason)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, m)
}

type batchCreateMessagesRequest struct {
	Messages []createMessageRequest `json:"messages"`
}

func (s *Server) batchCreateMessages(w http.ResponseWriter, r *http.Request) {
	var req batchCreateMessagesRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	inputs := make([]model.Message, len(req.Messages))
	for i, m := range req.Messages {
		inputs[i] = model.Message{TopicID: m.TopicID, Partition: m.Partition, Key: m.Key, Payload: m.Payload}
	}
	results, err := s.svc.BatchCreateMessages(inputs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, results)
}

func (s *Server) getMessagesByTopic(w http.ResponseWriter, r *http.Request) {
	topicID := r.PathValue("topic_id")
	items, err := s.svc.GetMessagesByTopic(topicID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, items)
}
