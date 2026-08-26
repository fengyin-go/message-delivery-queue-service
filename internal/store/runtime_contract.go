package store

import (
	"errors"
)

type OptionalTopicPolicyValidator interface{ Validate(string) error }
type OptionalTopicPolicyRuleSet struct{ rules map[string]bool }

func (r *OptionalTopicPolicyRuleSet) Validate(key string) error {
	if !r.rules[key] {
		return errors.New("route rejected")
	}
	return nil
}
func LoadOptionalTopicPolicyValidator(enabled bool) OptionalTopicPolicyValidator {
	if !enabled {
		var empty *OptionalTopicPolicyRuleSet
		return empty
	}
	return &OptionalTopicPolicyRuleSet{}
}
func OptionalTopicPolicyValidatorUsable(v OptionalTopicPolicyValidator) bool { return v != nil }
