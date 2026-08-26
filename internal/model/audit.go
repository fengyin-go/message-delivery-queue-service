package model

import (
	"strings"
	"time"
)

const (
	AuditActionCreate = "create"
	AuditActionUpdate = "update"
	AuditActionDelete = "delete"
	AuditActionDeliver = "deliver"
	AuditActionAck = "ack"
	AuditActionFail = "fail"
)

// Audit 审计日志实体。
type Audit struct {
	ID         string    `json:"id"`
	EntityType string    `json:"entity_type"`
	EntityID   string    `json:"entity_id"`
	Action     string    `json:"action"`
	Operator   string    `json:"operator"`
	Detail     string    `json:"detail"`
	CreatedAt  time.Time `json:"created_at"`
}

func (a *Audit) Validate() error {
	a.EntityType = strings.TrimSpace(a.EntityType)
	a.Action = strings.TrimSpace(a.Action)
	if a.EntityType == "" {
		return NewValidationError("entity_type", "实体类型不能为空")
	}
	if a.EntityID == "" {
		return NewValidationError("entity_id", "实体 ID 不能为空")
	}
	if a.Action == "" {
		return NewValidationError("action", "操作类型不能为空")
	}
	validActions := map[string]bool{
		AuditActionCreate: true, AuditActionUpdate: true, AuditActionDelete: true,
		AuditActionDeliver: true, AuditActionAck: true, AuditActionFail: true,
	}
	if !validActions[a.Action] {
		return NewValidationError("action", "操作类型不合法")
	}
	return nil
}

type AuditFilter struct {
	EntityType string
	EntityID   string
	Action     string
	Operator   string
}

func (f AuditFilter) Match(a *Audit) bool {
	if f.EntityType != "" && a.EntityType != f.EntityType {
		return false
	}
	if f.EntityID != "" && a.EntityID != f.EntityID {
		return false
	}
	if f.Action != "" && a.Action != f.Action {
		return false
	}
	if f.Operator != "" && a.Operator != f.Operator {
		return false
	}
	return true
}
