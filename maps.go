package govalid

import (
	"errors"
	"reflect"
)

func Map() Rule { return mapRule(func(reflect.Value) error { return nil }) }

// MapIsMap is kept for backward compatibility. Use Map instead.
func MapIsMap() Rule { return Map() }

func MapNotNil() Rule {
	return mapRule(func(value reflect.Value) error {
		if value.IsNil() {
			return errors.New("should not be nil")
		}
		return nil
	})
}

func MapNil() Rule {
	return mapRule(func(value reflect.Value) error {
		if !value.IsNil() {
			return errors.New("should be nil")
		}
		return nil
	})
}

func MapEmpty() Rule {
	return mapRule(func(value reflect.Value) error {
		if value.Len() != 0 {
			return errors.New("should be an empty map")
		}
		return nil
	})
}

func MapNotEmpty() Rule {
	return mapRule(func(value reflect.Value) error {
		if value.Len() == 0 {
			return errors.New("should be a non-empty map")
		}
		return nil
	})
}
