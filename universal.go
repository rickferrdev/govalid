package govalid

import (
	"errors"
	"reflect"
)

// Required rejects nil values and empty strings, arrays, slices, and maps.
// Scalar zero values such as 0 and false are considered present; combine it
// with NotZero when a non-zero scalar is required.
func Required() Rule {
	return func(context RuleContext) error {
		if isAbsentValue(context.Value) {
			return errors.New("should be present and non-empty")
		}
		return nil
	}
}

// Nil requires a nil-capable value to be nil.
func Nil() Rule {
	return func(context RuleContext) error {
		if !isNilValue(context.Value) {
			return errors.New("should be nil")
		}
		return nil
	}
}

// NotNil rejects nil values. Values whose kinds cannot be nil always pass.
func NotNil() Rule {
	return func(context RuleContext) error {
		if isNilValue(context.Value) {
			return errors.New("should not be nil")
		}
		return nil
	}
}

// Zero requires the field to be the zero value of its type.
func Zero() Rule {
	return func(context RuleContext) error {
		if !context.Value.IsValid() || !context.Value.IsZero() {
			return errors.New("should be the zero value")
		}
		return nil
	}
}

// NotZero rejects the zero value of the field type.
func NotZero() Rule {
	return func(context RuleContext) error {
		if !context.Value.IsValid() || context.Value.IsZero() {
			return errors.New("should not be the zero value")
		}
		return nil
	}
}

func isAbsentValue(value reflect.Value) bool {
	if isNilValue(value) {
		return true
	}

	value = indirect(value)
	if !value.IsValid() {
		return true
	}

	switch value.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return value.Len() == 0
	default:
		return false
	}
}
