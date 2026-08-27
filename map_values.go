package govalid

import (
	"errors"
	"fmt"
	"reflect"
)

func MapContainsValue[V any](expected V) Rule {
	return mapRule(func(value reflect.Value) error {
		iterator := value.MapRange()
		for iterator.Next() {
			if reflect.DeepEqual(iterator.Value().Interface(), expected) {
				return nil
			}
		}
		return fmt.Errorf("map should contain value %v", expected)
	})
}

func MapNotContainsValue[V any](unexpected V) Rule {
	return mapRule(func(value reflect.Value) error {
		iterator := value.MapRange()
		for iterator.Next() {
			if reflect.DeepEqual(iterator.Value().Interface(), unexpected) {
				return fmt.Errorf("map should not contain value %v", unexpected)
			}
		}
		return nil
	})
}

func MapContainsAnyValue[V any](expected ...V) Rule {
	return mapRule(func(value reflect.Value) error {
		iterator := value.MapRange()
		for iterator.Next() {
			for _, expectedValue := range expected {
				if reflect.DeepEqual(iterator.Value().Interface(), expectedValue) {
					return nil
				}
			}
		}
		return fmt.Errorf("map should contain at least one of the values %v", expected)
	})
}

func MapContainsAllValues[V any](expected ...V) Rule {
	return mapRule(func(value reflect.Value) error {
		for _, expectedValue := range expected {
			found := false
			iterator := value.MapRange()
			for iterator.Next() {
				if reflect.DeepEqual(iterator.Value().Interface(), expectedValue) {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("map should contain value %v", expectedValue)
			}
		}
		return nil
	})
}

func MapValues(rules ...Rule) Rule {
	return func(context ruleContext) error {
		value, err := mapValue(context)
		if err != nil {
			return err
		}

		iterator := value.MapRange()
		for iterator.Next() {
			key := iterator.Key()
			valueContext := ruleContext{
				Root:  context.Root,
				Path:  mapEntryPath(context.Path, key),
				Value: indirect(iterator.Value()),
			}
			if err := applyRules(valueContext, rules); err != nil {
				return fmt.Errorf("invalid map value at %s: %w", valueContext.Path, err)
			}
		}
		return nil
	}
}

func MapValueAt[K comparable](key K, rules ...Rule) Rule {
	return mapValueAt(key, false, rules)
}

func MapValueAtIfPresent[K comparable](key K, rules ...Rule) Rule {
	return mapValueAt(key, true, rules)
}

func MapNoNilValues() Rule {
	return mapRule(func(value reflect.Value) error {
		iterator := value.MapRange()
		for iterator.Next() {
			if isNilValue(iterator.Value()) {
				return fmt.Errorf("map value at key %v should not be nil", iterator.Key().Interface())
			}
		}
		return nil
	})
}

func MapNoZeroValues() Rule {
	return mapRule(func(value reflect.Value) error {
		iterator := value.MapRange()
		for iterator.Next() {
			item := indirect(iterator.Value())
			if !item.IsValid() || item.IsZero() {
				return fmt.Errorf("map value at key %v should not be zero", iterator.Key().Interface())
			}
		}
		return nil
	})
}

func MapHasNonZeroValue() Rule {
	return mapRule(func(value reflect.Value) error {
		iterator := value.MapRange()
		for iterator.Next() {
			item := indirect(iterator.Value())
			if item.IsValid() && !item.IsZero() {
				return nil
			}
		}
		return errors.New("map should contain at least one non-zero value")
	})
}

func mapValueAt[K comparable](expectedKey K, optional bool, rules []Rule) Rule {
	return func(context ruleContext) error {
		value, err := mapValue(context)
		if err != nil {
			return err
		}
		key, err := mapKeyValue(value, expectedKey)
		if err != nil {
			return err
		}
		item := value.MapIndex(key)
		if !item.IsValid() {
			if optional {
				return nil
			}
			return fmt.Errorf("map should contain key %v", expectedKey)
		}
		itemContext := ruleContext{
			Root: context.Root, Path: mapEntryPath(context.Path, key), Value: indirect(item),
		}
		if err := applyRules(itemContext, rules); err != nil {
			return fmt.Errorf("invalid map value at %s: %w", itemContext.Path, err)
		}
		return nil
	}
}
