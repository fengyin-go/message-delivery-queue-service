package handler

import (
	"net/http"
	"strconv"

	"message-queue/pkg/httpx"
)

func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats/global", s.getGlobalStats)
	mux.HandleFunc("GET /api/stats/topics", s.getTopicMessageStats)
	mux.HandleFunc("GET /api/stats/rates", s.getTopicRateStats)
	mux.HandleFunc("GET /api/stats/dead-letter-topn", s.getDeadLetterTopN)
	mux.HandleFunc("GET /api/topics/{topic_id}/export", s.exportTopic)
}

func (s *Server) getGlobalStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.svc.GetGlobalStats()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, stats)
}

func (s *Server) getTopicMessageStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.svc.GetTopicMessageStats()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, stats)
}

func (s *Server) getTopicRateStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.svc.GetTopicRateStats()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, stats)
}

func (s *Server) getDeadLetterTopN(w http.ResponseWriter, r *http.Request) {
	nStr := r.URL.Query().Get("n")
	n, _ := strconv.Atoi(nStr)
	if n <= 0 {
		n = 5
	}
	stats, err := s.svc.GetDeadLetterTopN(n)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, stats)
}

func (s *Server) exportTopic(w http.ResponseWriter, r *http.Request) {
	topicID := r.PathValue("topic_id")
	export, err := s.svc.ExportTopic(topicID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, export)
}
