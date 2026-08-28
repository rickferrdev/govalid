package govalid

import (
	"errors"
	"reflect"
)

// CollectionEqual returns a collection validation rule for equal.
func CollectionEqual[T any](expected T) Rule {
	return collectionRule(func(value reflect.Value) error {
		if !reflect.DeepEqual(value.Interface(), expected) {
			return errors.New("collection should be equal to the expected collection")
		}
		return nil
	})
}

// CollectionNotEqual returns a collection validation rule for not equal.
func CollectionNotEqual[T any](unexpected T) Rule {
	return collectionRule(func(value reflect.Value) error {
		if reflect.DeepEqual(value.Interface(), unexpected) {
			return errors.New("collection should be different from the unexpected collection")
		}
		return nil
	})
}
