package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// INT_MIN.
type intMinRule struct {
	typeRule TypeValidationRule
	min      int64
}

func (rule intMinRule) Validate(fieldName string, fieldValue reflect.Value) ValidationErrors {
	return rule.typeRule.Validate(fieldName, fieldValue, func(name string, value reflect.Value) *ValidationError {
		if value.Int() < rule.min {
			return &ValidationError{name, fmt.Errorf("%w: field value is lower than min", ErrIntMinRuleIsInvalid)}
		}

		return nil
	})
}

func NewIntMinRule(ruleStr string) (ValidationRule, error) {
	ruleValue := parseRule(ruleStr, "min:")
	v, err := strconv.ParseInt(ruleValue, 10, 0)
	if err != nil {
		return nil, ErrIntMinRuleWrongFormat
	}

	return intMinRule{typeRule{reflect.Int}, v}, nil
}

// INT_MAX.
type intMaxRule struct {
	typeRule TypeValidationRule
	max      int64
}

func (rule intMaxRule) Validate(fieldName string, fieldValue reflect.Value) ValidationErrors {
	return rule.typeRule.Validate(fieldName, fieldValue, func(name string, value reflect.Value) *ValidationError {
		if value.Int() > rule.max {
			return &ValidationError{name, fmt.Errorf("%w: field value is bigger than max", ErrIntMaxRuleIsInvalid)}
		}

		return nil
	})
}

func NewIntMaxRule(ruleStr string) (ValidationRule, error) {
	ruleValue := parseRule(ruleStr, "max:")
	v, err := strconv.ParseInt(ruleValue, 10, 0)
	if err != nil {
		return nil, ErrIntMaxRuleWrongFormat
	}

	return intMaxRule{typeRule{reflect.Int}, v}, nil
}

// INT_IN.
type intInRule struct {
	typeRule            TypeValidationRule
	availableValuesList []int64
}

func (rule intInRule) Validate(fieldName string, fieldValue reflect.Value) ValidationErrors {
	return rule.typeRule.Validate(fieldName, fieldValue, func(name string, value reflect.Value) *ValidationError {
		for _, availableValue := range rule.availableValuesList {
			if availableValue == value.Int() {
				return nil
			}
		}

		return &ValidationError{
			name,
			fmt.Errorf("%w: field value is not matching any of %v", ErrIntInRuleIsInvalid, rule.availableValuesList),
		}
	})
}

func NewIntInRule(ruleStr string) (ValidationRule, error) {
	ruleValue := parseRule(ruleStr, "in:")
	strValues := strings.Split(ruleValue, ValidateValuesSep)
	valuesList := make([]int64, 0, len(strValues))
	for _, strValue := range strValues {
		parsedInt, err := strconv.ParseInt(strValue, 10, 0)
		if err != nil {
			return nil, ErrIntInRuleWrongFormat
		}

		valuesList = append(valuesList, parsedInt)
	}

	if len(valuesList) == 0 {
		return nil, ErrIntInRuleWrongFormat
	}

	return intInRule{typeRule{reflect.Int}, valuesList}, nil
}
