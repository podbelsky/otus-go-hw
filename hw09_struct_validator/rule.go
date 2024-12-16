package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strings"
)

type typeRule struct {
	kind reflect.Kind
}

func (rule typeRule) Validate(fieldName string, fieldValue reflect.Value, logic validationLogic) ValidationErrors {
	switch fieldValue.Kind() { //nolint:exhaustive
	case rule.kind:
		if err := logic(fieldName, fieldValue); err != nil {
			return []ValidationError{*err}
		}

		return nil
	case reflect.Slice, reflect.Array:
		errors := ValidationErrors{}

		if fieldValue.IsNil() {
			return errors
		}

		for i := 0; i < fieldValue.Len(); i++ {
			elem := fieldValue.Index(i)
			if elem.Kind() != rule.kind {
				errors = append(errors, ValidationError{
					fieldName,
					fmt.Errorf("%w: one of field values type is not a %s type", ErrTypeRuleIsInvalid, rule.kind.String()),
				})
			}

			if err := logic(fieldName, elem); err != nil {
				errors = append(errors, *err)
			}
		}

		return errors
	default:
		return []ValidationError{{
			fieldName,
			fmt.Errorf("%w: field is not a %s or %s array", ErrTypeRuleIsInvalid, rule.kind.String(), rule.kind.String()),
		}}
	}
}

func parseRule(rule string, rulePrefix string) string {
	if len(rule) == 0 {
		return ""
	}

	if strings.HasPrefix(rule, rulePrefix) {
		ruleStrParts := strings.SplitN(rule, ":", 2)
		if len(ruleStrParts) > 1 {
			return ruleStrParts[1]
		}
	}

	return ""
}
