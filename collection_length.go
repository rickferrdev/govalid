package govalid

import (
	"errors"
	"fmt"
	"reflect"
)

func CollectionLength(expect int) Rule {
	return collectionLengthRule(expect,
		func(length int) bool { return length == expect },
		fmt.Sprintf("should have exactly %d items", expect))
}

func CollectionMinLength(expect int) Rule {
	return collectionLengthRule(expect,
		func(length int) bool { return length >= expect },
		fmt.Sprintf("should have at least %d items", expect))
}

func CollectionMaxLength(expect int) Rule {
	return collectionLengthRule(expect,
		func(length int) bool { return length <= expect },
		fmt.Sprintf("should have at most %d items", expect))
}

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
