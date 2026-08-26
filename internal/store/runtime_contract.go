package store

import (
	"errors"
)

type OptionalTopicPolicyValidator interface{ Validate(string) error }
type OptionalTopicPolicyRuleSet struct{ rules map[string]bool }

// Allow 登记一个主题，使其在校验时被放行。仅在策略启用后登记才有意义；
// 未登记的主题将在 Validate 时被拒绝。
func (r *OptionalTopicPolicyRuleSet) Allow(key string) {
	if r == nil || key == "" {
		return
	}
	if r.rules == nil {
		r.rules = make(map[string]bool)
	}
	r.rules[key] = true
}

func (r *OptionalTopicPolicyRuleSet) Validate(key string) error {
	if r == nil || r.rules == nil {
		// 策略启用但规则集为空：没有主题被登记，任何 key 都应被拒绝。
		return errors.New("route rejected")
	}
	if !r.rules[key] {
		return errors.New("route rejected")
	}
	return nil
}

// LoadOptionalTopicPolicyValidator 按 enabled 决定是否启用可选主题策略。
// 未启用时返回 nil 接口（而非包裹 nil 指针的接口），避免
// OptionalTopicPolicyValidatorUsable 误判为可用而进入空指针解引用。
func LoadOptionalTopicPolicyValidator(enabled bool) OptionalTopicPolicyValidator {
	if !enabled {
		return nil
	}
	return &OptionalTopicPolicyRuleSet{rules: make(map[string]bool)}
}

// OptionalTopicPolicyValidatorUsable 判断校验器是否真正可用。
// 同时防御未启用时返回的 nil 接口与启用后仍可能出现的 nil 规则集。
func OptionalTopicPolicyValidatorUsable(v OptionalTopicPolicyValidator) bool {
	if v == nil {
		return false
	}
	// 接口可能包裹 nil 指针（历史遗留），此时调用 Validate 会 panic，
	// 一律视为不可用，交由上层放行。
	if rs, ok := v.(*OptionalTopicPolicyRuleSet); ok && rs == nil {
		return false
	}
	return true
}
