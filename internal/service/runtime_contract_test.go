package service_test

import (
	"message-queue/internal/service"
	"message-queue/internal/store"
	"testing"
)

func callOptionalTopicPolicy(v store.OptionalTopicPolicyValidator, key string) (panicked bool, err error) {
	defer func() {
		if recover() != nil {
			panicked = true
		}
	}()
	err = service.ValidateOptionalTopicPolicy(v, key)
	return
}

func TestR022Target(t *testing.T) {
	disabledPanic, disabledErr := callOptionalTopicPolicy(store.LoadOptionalTopicPolicyValidator(false), "anything")
	enabledPanic, enabledErr := callOptionalTopicPolicy(store.LoadOptionalTopicPolicyValidator(true), "missing")
	if disabledPanic || disabledErr != nil || enabledPanic || enabledErr == nil {
		t.Fatalf("关闭规则时发生崩溃或启用规则后缺失值被放行")
	}
}
