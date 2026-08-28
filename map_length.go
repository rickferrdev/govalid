package govalid

import (
	"errors"
	"fmt"
	"reflect"
)

// MapLength returns a map validation rule for length.
func MapLength(expect int) Rule {
	return mapLengthRule(expect, func(length int) bool { return length == expect },
		fmt.Sprintf("should have exactly %d entries", expect))
}

// MapMinLength returns a map validation rule for min length.
func MapMinLength(expect int) Rule {
	return mapLengthRule(expect, func(length int) bool { return length >= expect },
		fmt.Sprintf("should have at least %d entries", expect))
}

// MapMaxLength returns a map validation rule for max length.
func MapMaxLength(expect int) Rule {
	return mapLengthRule(expect, func(length int) bool { return length <= expect },
		fmt.Sprintf("should have at most %d entries", expect))
}

// MapLengthBetween returns a map validation rule for length between.
func MapLengthBetween(minimum, maximum int) Rule {
	return mapRule(func(value reflect.Value) error {
		if minimum < 0 || maximum < 0 {
			return errors.New("map length cannot be negative")
		}
		if minimum > maximum {
			return errors.New("minimum cannot be greater than maximum")
		}
		if value.Len() < minimum || value.Len() > maximum {
			return fmt.Errorf("map length should be between %d and %d", minimum, maximum)
		}
		return nil
	})
}

// MapLengthNotBetween returns a map validation rule for length not between.
func MapLengthNotBetween(minimum, maximum int) Rule {
	return mapRule(func(value reflect.Value) error {
		if minimum < 0 || maximum < 0 {
			return errors.New("map length cannot be negative")
		}
		if minimum > maximum {
			return errors.New("minimum cannot be greater than maximum")
		}
		if value.Len() >= minimum && value.Len() <= maximum {
			return fmt.Errorf("map length should not be between %d and %d", minimum, maximum)
		}
		return nil
	})
}

func mapLengthRule(expect int, valid func(int) bool, message string) Rule {
	return mapRule(func(value reflect.Value) error {
		if expect < 0 {
			return errors.New("map length cannot be negative")
		}
		if !valid(value.Len()) {
			return errors.New(message)
		}
		return nil
	})
}
