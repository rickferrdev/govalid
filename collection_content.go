package govalid

import (
	"errors"
	"fmt"
	"reflect"
)

func CollectionContains[T any](expected T) Rule {
	return collectionRule(func(value reflect.Value) error {
		if collectionContains(value, expected) {
			return nil
		}
		return fmt.Errorf("collection should contain %v", expected)
	})
}

func CollectionNotContains[T any](unexpected T) Rule {
	return collectionRule(func(value reflect.Value) error {
		if collectionContains(value, unexpected) {
			return fmt.Errorf("collection should not contain %v", unexpected)
		}
		return nil
	})
}

func CollectionContainsAny[T any](expected ...T) Rule {
	return collectionRule(func(value reflect.Value) error {
		for _, expectedItem := range expected {
			if collectionContains(value, expectedItem) {
				return nil
			}
		}
		return fmt.Errorf("collection should contain at least one of %v", expected)
	})
}

func CollectionContainsAll[T any](expected ...T) Rule {
	return collectionRule(func(value reflect.Value) error {
		for _, expectedItem := range expected {
			if !collectionContains(value, expectedItem) {
				return fmt.Errorf("collection should contain %v", expectedItem)
			}
		}
		return nil
	})
}

func CollectionUnique() Rule {
	return collectionRule(func(value reflect.Value) error {
		for left := 0; left < value.Len(); left++ {
			for right := left + 1; right < value.Len(); right++ {
				if reflect.DeepEqual(value.Index(left).Interface(), value.Index(right).Interface()) {
					return fmt.Errorf("collection should contain unique items; duplicate at indexes %d and %d", left, right)
				}
			}
		}
		return nil
	})
}

func CollectionNoNilItems() Rule {
	return collectionRule(func(value reflect.Value) error {
		for index := 0; index < value.Len(); index++ {
			if isNilValue(value.Index(index)) {
				return fmt.Errorf("collection item at index %d should not be nil", index)
			}
		}
		return nil
	})
}

func CollectionNoZeroItems() Rule {
	return collectionRule(func(value reflect.Value) error {
		for index := 0; index < value.Len(); index++ {
			item := indirect(value.Index(index))
			if !item.IsValid() || item.IsZero() {
				return fmt.Errorf("collection item at index %d should not be zero", index)
			}
		}
		return nil
	})
}

func CollectionHasNonZeroItem() Rule {
	return collectionRule(func(value reflect.Value) error {
		for index := 0; index < value.Len(); index++ {
			item := indirect(value.Index(index))
			if item.IsValid() && !item.IsZero() {
				return nil
			}
		}
		return errors.New("collection should contain at least one non-zero item")
	})
}

func collectionContains[T any](value reflect.Value, expected T) bool {
	for index := 0; index < value.Len(); index++ {
		if reflect.DeepEqual(value.Index(index).Interface(), expected) {
			return true
		}
	}
	return false
}
