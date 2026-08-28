package govalid

import (
	"errors"
	"reflect"
)

// Map requires a map value.
func Map() Rule { return mapRule(func(reflect.Value) error { return nil }) }

// MapNotNil returns a map validation rule for not nil.
func MapNotNil() Rule {
	return mapRule(func(value reflect.Value) error {
		if value.IsNil() {
			return errors.New("should not be nil")
		}
		return nil
	})
}

// MapNil returns a map validation rule for nil.
func MapNil() Rule {
	return mapRule(func(value reflect.Value) error {
		if !value.IsNil() {
			return errors.New("should be nil")
		}
		return nil
	})
}

// MapEmpty returns a map validation rule for empty.
func MapEmpty() Rule {
	return mapRule(func(value reflect.Value) error {
		if value.Len() != 0 {
			return errors.New("should be an empty map")
		}
		return nil
	})
}

// MapNotEmpty returns a map validation rule for not empty.
func MapNotEmpty() Rule {
	return mapRule(func(value reflect.Value) error {
		if value.Len() == 0 {
			return errors.New("should be a non-empty map")
		}
		return nil
	})
}
