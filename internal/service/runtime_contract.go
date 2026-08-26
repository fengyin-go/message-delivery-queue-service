package service

import "message-queue/internal/store"

func ValidateOptionalTopicPolicy(v store.OptionalTopicPolicyValidator, key string) (err error) {

	if !store.OptionalTopicPolicyValidatorUsable(v) {
		return nil
	}
	return v.Validate(key)
}
