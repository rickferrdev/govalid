package govalid

import (
	"errors"
	"reflect"
)

// Bytes requires a byte slice or byte array, including defined byte types.
func Bytes() Rule {
	return bytesRule(func(reflect.Value) error { return nil })
}

// BytesNil returns a byte-sequence validation rule for nil.
func BytesNil() Rule {
	return bytesRule(func(value reflect.Value) error {
		if value.Kind() != reflect.Slice || !value.IsNil() {
			return errors.New("should be nil bytes")
		}
		return nil
	})
}

// BytesNotNil returns a byte-sequence validation rule for not nil.
func BytesNotNil() Rule {
	return bytesRule(func(value reflect.Value) error {
		if value.Kind() == reflect.Slice && value.IsNil() {
			return errors.New("should be non-nil bytes")
		}
		return nil
	})
}

// BytesEmpty returns a byte-sequence validation rule for empty.
func BytesEmpty() Rule {
	return bytesPredicateRule(
		func(value reflect.Value) bool { return value.Len() == 0 },
		"should be empty bytes",
	)
}

// BytesNotEmpty returns a byte-sequence validation rule for not empty.
func BytesNotEmpty() Rule {
	return bytesPredicateRule(
		func(value reflect.Value) bool { return value.Len() > 0 },
		"should be non-empty bytes",
	)
}
