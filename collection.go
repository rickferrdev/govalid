package govalid

import (
	"errors"
	"reflect"
)

// Collection requires a slice or array value.
func Collection() Rule {
	return collectionRule(func(reflect.Value) error { return nil })
}

// CollectionNotNil returns a collection validation rule for not nil.
func CollectionNotNil() Rule {
	return collectionRule(func(value reflect.Value) error {
		if value.Kind() == reflect.Slice && value.IsNil() {
			return errors.New("should be a non-nil collection")
		}
		return nil
	})
}

// CollectionNil returns a collection validation rule for nil.
func CollectionNil() Rule {
	return collectionRule(func(value reflect.Value) error {
		if value.Kind() != reflect.Slice || !value.IsNil() {
			return errors.New("should be a nil collection")
		}
		return nil
	})
}

// CollectionEmpty returns a collection validation rule for empty.
func CollectionEmpty() Rule {
	return collectionPredicateRule(
		func(value reflect.Value) bool { return value.Len() == 0 },
		"should be an empty collection",
	)
}

// CollectionNotEmpty returns a collection validation rule for not empty.
func CollectionNotEmpty() Rule {
	return collectionPredicateRule(
		func(value reflect.Value) bool { return value.Len() > 0 },
		"should be a non-empty collection",
	)
}
