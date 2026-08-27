package govalid

import (
	"errors"
	"reflect"
)

func Bytes() Rule {
	return bytesRule(func(reflect.Value) error { return nil })
}

func BytesNil() Rule {
	return bytesRule(func(value reflect.Value) error {
		if value.Kind() != reflect.Slice || !value.IsNil() {
			return errors.New("should be nil bytes")
		}
		return nil
	})
}

func BytesNotNil() Rule {
	return bytesRule(func(value reflect.Value) error {
		if value.Kind() == reflect.Slice && value.IsNil() {
			return errors.New("should be non-nil bytes")
		}
		return nil
	})
}

func BytesEmpty() Rule {
	return bytesPredicateRule(
		func(value reflect.Value) bool { return value.Len() == 0 },
		"should be empty bytes",
	)
}

func BytesNotEmpty() Rule {
	return bytesPredicateRule(
		func(value reflect.Value) bool { return value.Len() > 0 },
		"should be non-empty bytes",
	)
}
