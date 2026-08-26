package store

import (
	"errors"
)

type OptionalAuditFilterValidator interface{ Validate(string) error }
type OptionalAuditFilterRuleSet struct{ rules map[string]bool }

func (r *OptionalAuditFilterRuleSet) Validate(key string) error {
	if !r.rules[key] {
		return errors.New("route rejected")
	}
	return nil
}
func LoadOptionalAuditFilterValidator(enabled bool) OptionalAuditFilterValidator {
	if !enabled {
		var empty *OptionalAuditFilterRuleSet
		return empty
	}
	return &OptionalAuditFilterRuleSet{}
}
func OptionalAuditFilterValidatorUsable(v OptionalAuditFilterValidator) bool { return v != nil }
