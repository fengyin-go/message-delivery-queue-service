package service

import "message-queue/internal/store"

func ValidateOptionalRoutePolicy(v store.OptionalRoutePolicyValidator, key string) (err error) {

	if !store.OptionalRoutePolicyValidatorUsable(v) {
		return nil
	}
	return v.Validate(key)
}
