package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// STRING_LEN.
type stringLenRule struct {
	typeRule TypeValidationRule
	len      int
}

func (rule stringLenRule) Validate(fieldName string, fieldValue reflect.Value) ValidationErrors {
	return rule.typeRule.Validate(fieldName, fieldValue, func(name string, value reflect.Value) *ValidationError {
		if len(value.String()) != rule.len {
			return &ValidationError{name, fmt.Errorf("%w: field length not equals %d", ErrStrLengthRuleIsInvalid, rule.len)}
		}

		return nil
	})
}

func NewStringLenRule(ruleStr string) (ValidationRule, error) {
	ruleValue := parseRule(ruleStr, "len:")
	length, err := strconv.Atoi(ruleValue)
	if err != nil {
		return nil, ErrStrLengthRuleWrongFormat
	}

	return stringLenRule{typeRule{reflect.String}, length}, nil
}

// STRING_IN.
type stringInRule struct {
	typeRule            TypeValidationRule
	availableValuesList []string
}

func NewStringInRule(ruleStr string) (ValidationRule, error) {
	ruleValue := parseRule(ruleStr, "in:")
	valuesList := strings.Split(ruleValue, ValidateValuesSep)
	if len(valuesList) == 0 {
		return nil, ErrStrInRuleWrongFormat
	}
	return stringInRule{typeRule{reflect.String}, valuesList}, nil
}

func (rule stringInRule) Validate(fieldName string, fieldValue reflect.Value) ValidationErrors {
	return rule.typeRule.Validate(fieldName, fieldValue, func(name string, value reflect.Value) *ValidationError {
		for _, availableStr := range rule.availableValuesList {
			if availableStr == value.String() {
				return nil
			}
		}
		return &ValidationError{
			name,
			fmt.Errorf("%w field value is not matching any of %v", ErrStrInRuleIsInvalid, rule.availableValuesList),
		}
	})
}

// REGEXP..
type stringRegexpRule struct {
	typeRule TypeValidationRule
	regexp   *regexp.Regexp
}

func NewStringRegexpRule(ruleStr string) (ValidationRule, error) {
	ruleValue := parseRule(ruleStr, "regexp:")
	r, err := regexp.Compile(ruleValue)
	if err != nil {
		return nil, ErrStrRegexpRuleWrongFormat
	}

	return stringRegexpRule{typeRule{reflect.String}, r}, nil
}

func (rule stringRegexpRule) Validate(fieldName string, fieldValue reflect.Value) ValidationErrors {
	return rule.typeRule.Validate(fieldName, fieldValue, func(name string, value reflect.Value) *ValidationError {
		if !rule.regexp.MatchString(value.String()) {
			return &ValidationError{name, fmt.Errorf("%w: field is not matching regular expression", ErrStrRegexpRuleIsInvalid)}
		}
		return nil
	})
}
