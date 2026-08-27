package govalid

import (
	"fmt"
	"reflect"
)

func MapHasKey[K comparable](expected K) Rule {
	return mapRule(func(value reflect.Value) error {
		key, err := mapKeyValue(value, expected)
		if err != nil {
			return err
		}
		if !value.MapIndex(key).IsValid() {
			return fmt.Errorf("map should contain key %v", expected)
		}
		return nil
	})
}

func MapNotHasKey[K comparable](unexpected K) Rule {
	return mapRule(func(value reflect.Value) error {
		key, err := mapKeyValue(value, unexpected)
		if err != nil {
			return err
		}
		if value.MapIndex(key).IsValid() {
			return fmt.Errorf("map should not contain key %v", unexpected)
		}
		return nil
	})
}

func MapHasAllKeys[K comparable](expected ...K) Rule {
	return mapRule(func(value reflect.Value) error {
		for _, expectedKey := range expected {
			key, err := mapKeyValue(value, expectedKey)
			if err != nil {
				return err
			}
			if !value.MapIndex(key).IsValid() {
				return fmt.Errorf("map should contain key %v", expectedKey)
			}
		}
		return nil
	})
}

func MapHasAnyKey[K comparable](expected ...K) Rule {
	return mapRule(func(value reflect.Value) error {
		for _, expectedKey := range expected {
			key, err := mapKeyValue(value, expectedKey)
			if err != nil {
				return err
			}
			if value.MapIndex(key).IsValid() {
				return nil
			}
		}
		return fmt.Errorf("map should contain at least one of the keys %v", expected)
	})
}

func MapHasNoneOfKeys[K comparable](unexpected ...K) Rule {
	return mapRule(func(value reflect.Value) error {
		for _, unexpectedKey := range unexpected {
			key, err := mapKeyValue(value, unexpectedKey)
			if err != nil {
				return err
			}
			if value.MapIndex(key).IsValid() {
				return fmt.Errorf("map should not contain key %v", unexpectedKey)
			}
		}
		return nil
	})
}

func MapAllowedKeys[K comparable](allowed ...K) Rule {
	return mapRule(func(value reflect.Value) error {
		allowedKeys, err := mapKeyValues(value, allowed)
		if err != nil {
			return err
		}
		iterator := value.MapRange()
		for iterator.Next() {
			if !containsReflectValue(allowedKeys, iterator.Key()) {
				return fmt.Errorf("map contains an unexpected key %v", iterator.Key().Interface())
			}
		}
		return nil
	})
}

func MapHasExactKeys[K comparable](expected ...K) Rule {
	return mapRule(func(value reflect.Value) error {
		expectedKeys, err := mapKeyValues(value, expected)
		if err != nil {
			return err
		}
		if value.Len() != len(uniqueReflectValues(expectedKeys)) {
			return fmt.Errorf("map should contain exactly the keys %v", expected)
		}
		for _, key := range expectedKeys {
			if !value.MapIndex(key).IsValid() {
				return fmt.Errorf("map should contain exactly the keys %v", expected)
			}
		}
		return nil
	})
}

// Compatibility aliases for the previous API.
func MapHasAllKey[K comparable](expected ...K) Rule    { return MapHasAllKeys(expected...) }
func MapHasAnyAllKey[K comparable](expected ...K) Rule { return MapHasAnyKey(expected...) }
func MapHasNoneAllKey[K comparable](unexpected ...K) Rule {
	return MapHasNoneOfKeys(unexpected...)
}

func MapKeys(rules ...Rule) Rule {
	return func(context ruleContext) error {
		value, err := mapValue(context)
		if err != nil {
			return err
		}

		iterator := value.MapRange()
		for iterator.Next() {
			rawKey := iterator.Key()
			keyContext := ruleContext{
				Root:  context.Root,
				Path:  mapEntryPath(context.Path, rawKey),
				Value: indirect(rawKey),
			}
			if err := applyRules(keyContext, rules); err != nil {
				return fmt.Errorf("invalid map key at %s: %w", keyContext.Path, err)
			}
		}
		return nil
	}
}
