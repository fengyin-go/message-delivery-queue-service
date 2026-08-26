package store

import (
	"errors"
)

type OptionalRoutePolicyValidator interface{ Validate(string) error }
type OptionalRoutePolicyRuleSet struct{ rules map[string]bool }

func (r *OptionalRoutePolicyRuleSet) Validate(key string) error {
	if !r.rules[key] {
		return errors.New("route rejected")
	}
	return nil
}
func LoadOptionalRoutePolicyValidator(enabled bool) OptionalRoutePolicyValidator {
	if !enabled {
		var empty *OptionalRoutePolicyRuleSet
		return empty
	}
	return &OptionalRoutePolicyRuleSet{}
}
func OptionalRoutePolicyValidatorUsable(v OptionalRoutePolicyValidator) bool { return v != nil }
