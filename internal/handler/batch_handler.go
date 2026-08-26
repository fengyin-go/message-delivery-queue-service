package handler

import (
	"net/http"

	"message-queue/internal/model"
	"message-queue/pkg/httpx"
)

func (s *Server) registerBatchRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/batch/topics", s.batchCreateTopics)
	mux.HandleFunc("POST /api/batch/producers", s.batchCreateProducers)
	mux.HandleFunc("POST /api/batch/consumers", s.batchCreateConsumers)
	mux.HandleFunc("POST /api/batch/consumer-groups", s.batchCreateConsumerGroups)
	mux.HandleFunc("DELETE /api/batch/messages", s.batchDeleteMessages)
	mux.HandleFunc("DELETE /api/batch/dead-letters", s.batchDeleteDeadLetters)
}

type batchCreateTopicRequest struct {
	Topics []createTopicRequest `json:"topics"`
}

func (s *Server) batchCreateTopics(w http.ResponseWriter, r *http.Request) {
	var req batchCreateTopicRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	results := make([]*model.Topic, 0, len(req.Topics))
	for _, t := range req.Topics {
		created, err := s.svc.CreateTopic(model.Topic{Name: t.Name, Partitions: t.Partitions, RetentionSeconds: t.RetentionSeconds, Status: t.Status})
		if err != nil {
			writeServiceError(w, err)
			return
		}
		results = append(results, created)
	}
	httpx.Created(w, results)
}

type batchCreateProducerRequest struct {
	Producers []createProducerRequest `json:"producers"`
}

func (s *Server) batchCreateProducers(w http.ResponseWriter, r *http.Request) {
	var req batchCreateProducerRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	results := make([]*model.Producer, 0, len(req.Producers))
	for _, p := range req.Producers {
		created, err := s.svc.CreateProducer(model.Producer{Name: p.Name, TopicID: p.TopicID})
		if err != nil {
			writeServiceError(w, err)
			return
		}
		results = append(results, created)
	}
	httpx.Created(w, results)
}

type batchCreateConsumerRequest struct {
	Consumers []createConsumerRequest `json:"consumers"`
}

func (s *Server) batchCreateConsumers(w http.ResponseWriter, r *http.Request) {
	var req batchCreateConsumerRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	results := make([]*model.Consumer, 0, len(req.Consumers))
	for _, c := range req.Consumers {
		created, err := s.svc.CreateConsumer(model.Consumer{Name: c.Name})
		if err != nil {
			writeServiceError(w, err)
			return
		}
		results = append(results, created)
	}
	httpx.Created(w, results)
}

type batchCreateConsumerGroupRequest struct {
	Groups []createConsumerGroupRequest `json:"groups"`
}

func (s *Server) batchCreateConsumerGroups(w http.ResponseWriter, r *http.Request) {
	var req batchCreateConsumerGroupRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	results := make([]*model.ConsumerGroup, 0, len(req.Groups))
	for _, g := range req.Groups {
		created, err := s.svc.CreateConsumerGroup(model.ConsumerGroup{Name: g.Name, TopicID: g.TopicID})
		if err != nil {
			writeServiceError(w, err)
			return
		}
		results = append(results, created)
	}
	httpx.Created(w, results)
}

type batchDeleteMessagesRequest struct {
	IDs []string `json:"ids"`
}

func (s *Server) batchDeleteMessages(w http.ResponseWriter, r *http.Request) {
	var req batchDeleteMessagesRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	for _, id := range req.IDs {
		if err := s.svc.DeleteMessage(id); err != nil {
			writeServiceError(w, err)
			return
		}
	}
	httpx.NoContent(w)
}

type batchDeleteDeadLettersRequest struct {
	IDs []string `json:"ids"`
}

func (s *Server) batchDeleteDeadLetters(w http.ResponseWriter, r *http.Request) {
	var req batchDeleteDeadLettersRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	for _, id := range req.IDs {
		if err := s.svc.DeleteDeadLetter(id); err != nil {
			writeServiceError(w, err)
			return
		}
	}
	httpx.NoContent(w)
}
