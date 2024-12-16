package hw09structvalidator

import "reflect"

type validationLogic = func(fieldName string, fieldValue reflect.Value) *ValidationError

type ValidationRule interface {
	Validate(fieldName string, fieldValue reflect.Value) ValidationErrors
}

type TypeValidationRule interface {
	Validate(fieldName string, fieldValue reflect.Value, logic validationLogic) ValidationErrors
}
