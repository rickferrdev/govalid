package govalid

import (
	"errors"
	"fmt"
	"reflect"
)

// CollectionLength returns a collection validation rule for length.
func CollectionLength(expect int) Rule {
	return collectionLengthRule(expect,
		func(length int) bool { return length == expect },
		fmt.Sprintf("should have exactly %d items", expect))
}

// CollectionMinLength returns a collection validation rule for min length.
func CollectionMinLength(expect int) Rule {
	return collectionLengthRule(expect,
		func(length int) bool { return length >= expect },
		fmt.Sprintf("should have at least %d items", expect))
}

// CollectionMaxLength returns a collection validation rule for max length.
func CollectionMaxLength(expect int) Rule {
	return collectionLengthRule(expect,
		func(length int) bool { return length <= expect },
		fmt.Sprintf("should have at most %d items", expect))
}

// CollectionLengthBetween returns a collection validation rule for length between.
func CollectionLengthBetween(minimum, maximum int) Rule {
	return collectionRule(func(value reflect.Value) error {
		if minimum < 0 || maximum < 0 {
			return errors.New("collection length cannot be negative")
		}
		if minimum > maximum {
			return errors.New("minimum cannot be greater than maximum")
		}
		if value.Len() < minimum || value.Len() > maximum {
			return fmt.Errorf("collection length should be between %d and %d", minimum, maximum)
		}
		return nil
	})
}

// CollectionLengthNotBetween returns a collection validation rule for length not between.
func CollectionLengthNotBetween(minimum, maximum int) Rule {
	return collectionRule(func(value reflect.Value) error {
		if minimum < 0 || maximum < 0 {
			return errors.New("collection length cannot be negative")
		}
		if minimum > maximum {
			return errors.New("minimum cannot be greater than maximum")
		}
		if value.Len() >= minimum && value.Len() <= maximum {
			return fmt.Errorf("collection length should not be between %d and %d", minimum, maximum)
		}
		return nil
	})
}

func collectionLengthRule(expect int, valid func(int) bool, message string) Rule {
	return collectionRule(func(value reflect.Value) error {
		if expect < 0 {
			return errors.New("collection length cannot be negative")
		}
		if !valid(value.Len()) {
			return errors.New(message)
		}
		return nil
	})
}
