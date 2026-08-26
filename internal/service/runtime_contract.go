package service

import "message-queue/internal/store"

func ValidateOptionalAuditFilter(v store.OptionalAuditFilterValidator, key string) (err error) {

	if !store.OptionalAuditFilterValidatorUsable(v) {
		return nil
	}
	return v.Validate(key)
}
