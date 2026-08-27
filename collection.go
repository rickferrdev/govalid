package govalid

import (
	"errors"
	"reflect"
)

func Collection() Rule {
	return collectionRule(func(reflect.Value) error { return nil })
}

func CollectionNotNil() Rule {
	return collectionRule(func(value reflect.Value) error {
		if value.Kind() == reflect.Slice && value.IsNil() {
			return errors.New("should be a non-nil collection")
		}
		return nil
	})
}

func CollectionNil() Rule {
	return collectionRule(func(value reflect.Value) error {
		if value.Kind() != reflect.Slice || !value.IsNil() {
			return errors.New("should be a nil collection")
		}
		return nil
	})
}

func CollectionEmpty() Rule {
	return collectionPredicateRule(
		func(value reflect.Value) bool { return value.Len() == 0 },
		"should be an empty collection",
	)
}

func CollectionNotEmpty() Rule {
	return collectionPredicateRule(
		func(value reflect.Value) bool { return value.Len() > 0 },
		"should be a non-empty collection",
	)
}
