package govalid

import (
	"errors"
	"fmt"
	"reflect"
)

func mapValue(context ruleContext) (reflect.Value, error) {
	value := context.Value
	if !value.IsValid() || value.Kind() != reflect.Map {
		return reflect.Value{}, errors.New("should be of type map")
	}
	return value, nil
}

func mapRule(validate func(reflect.Value) error) Rule {
	return func(context ruleContext) error {
		value, err := mapValue(context)
		if err != nil {
			return err
		}
		return validate(value)
	}
}

func mapKeyValue[K comparable](value reflect.Value, key K) (reflect.Value, error) {
	keyValue := reflect.ValueOf(key)
	keyType := value.Type().Key()
	if !keyValue.IsValid() {
		if keyType.Kind() == reflect.Interface {
			return reflect.Zero(keyType), nil
		}
		return reflect.Value{}, fmt.Errorf("expected key type %s, received nil", keyType)
	}
	if !keyValue.Type().AssignableTo(keyType) {
		return reflect.Value{}, fmt.Errorf("expected key type %s, received %s", keyType, keyValue.Type())
	}
	if !keyValue.Comparable() {
		return reflect.Value{}, fmt.Errorf("key %v is not comparable", key)
	}
	return keyValue, nil
}

func mapKeyValues[K comparable](value reflect.Value, keys []K) ([]reflect.Value, error) {
	values := make([]reflect.Value, len(keys))
	for index, key := range keys {
		keyValue, err := mapKeyValue(value, key)
		if err != nil {
			return nil, err
		}
		values[index] = keyValue
	}
	return values, nil
}

func containsReflectValue(values []reflect.Value, expected reflect.Value) bool {
	for _, value := range values {
		if reflect.DeepEqual(value.Interface(), expected.Interface()) {
			return true
		}
	}
	return false
}

func uniqueReflectValues(values []reflect.Value) []reflect.Value {
	unique := make([]reflect.Value, 0, len(values))
	for _, value := range values {
		if !containsReflectValue(unique, value) {
			unique = append(unique, value)
		}
	}
	return unique
}

func mapEntryPath(path string, key reflect.Value) string {
	return fmt.Sprintf("%s[%v]", path, key.Interface())
}

func isNilValue(value reflect.Value) bool {
	for value.IsValid() && value.Kind() == reflect.Interface {
		if value.IsNil() {
			return true
		}
		value = value.Elem()
	}
	if !value.IsValid() {
		return true
	}
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}
