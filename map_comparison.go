package govalid

import (
	"errors"
	"reflect"
)

func MapEqual[M ~map[K]V, K comparable, V any](expected M) Rule {
	return mapRule(func(value reflect.Value) error {
		if !reflect.DeepEqual(value.Interface(), expected) {
			return errors.New("map should be equal to the expected map")
		}
		return nil
	})
}

func MapNotEqual[M ~map[K]V, K comparable, V any](unexpected M) Rule {
	return mapRule(func(value reflect.Value) error {
		if reflect.DeepEqual(value.Interface(), unexpected) {
			return errors.New("map should be different from the unexpected map")
		}
		return nil
	})
}

func MapSubsetOf[M ~map[K]V, K comparable, V any](expected M) Rule {
	return mapRule(func(value reflect.Value) error {
		expectedValue := reflect.ValueOf(expected)
		if value.Type() != expectedValue.Type() {
			return errors.New("expected map should have the same type as the validated map")
		}
		iterator := value.MapRange()
		for iterator.Next() {
			expectedItem := expectedValue.MapIndex(iterator.Key())
			if !expectedItem.IsValid() || !reflect.DeepEqual(iterator.Value().Interface(), expectedItem.Interface()) {
				return errors.New("map should be a subset of the expected map")
			}
		}
		return nil
	})
}

func MapSupersetOf[M ~map[K]V, K comparable, V any](expected M) Rule {
	return mapRule(func(value reflect.Value) error {
		expectedValue := reflect.ValueOf(expected)
		if value.Type() != expectedValue.Type() {
			return errors.New("expected map should have the same type as the validated map")
		}
		iterator := expectedValue.MapRange()
		for iterator.Next() {
			actualItem := value.MapIndex(iterator.Key())
			if !actualItem.IsValid() || !reflect.DeepEqual(actualItem.Interface(), iterator.Value().Interface()) {
				return errors.New("map should be a superset of the expected map")
			}
		}
		return nil
	})
}
