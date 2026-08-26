package handler

import (
	"net/http"

	"message-queue/internal/model"
	"message-queue/pkg/httpx"
)

func (s *Server) registerAuditRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/audits", s.createAudit)
	mux.HandleFunc("GET /api/audits", s.listAudits)
	mux.HandleFunc("GET /api/audits/{id}", s.getAudit)
	mux.HandleFunc("DELETE /api/audits/{id}", s.deleteAudit)
	mux.HandleFunc("GET /api/audits/entity/{entity_type}/{entity_id}", s.getAuditsByEntity)
}

type createAuditRequest struct {
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
	Action     string `json:"action"`
	Operator   string `json:"operator"`
	Detail     string `json:"detail"`
}

func (s *Server) createAudit(w http.ResponseWriter, r *http.Request) {
	var req createAuditRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	a, err := s.svc.CreateAudit(model.Audit{EntityType: req.EntityType, EntityID: req.EntityID, Action: req.Action, Operator: req.Operator, Detail: req.Detail})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, a)
}

func (s *Server) listAudits(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.AuditFilter{
		EntityType: r.URL.Query().Get("entity_type"),
		EntityID:   r.URL.Query().Get("entity_id"),
		Action:     r.URL.Query().Get("action"),
		Operator:   r.URL.Query().Get("operator"),
	}
	items, total, err := s.svc.ListAudits(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getAudit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	a, err := s.svc.GetAudit(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

func (s *Server) deleteAudit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteAudit(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) getAuditsByEntity(w http.ResponseWriter, r *http.Request) {
	entityType := r.PathValue("entity_type")
	entityID := r.PathValue("entity_id")
	items, err := s.svc.GetAuditsByEntity(entityType, entityID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, items)
}
