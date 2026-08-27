package govalid

import (
	"errors"
	"reflect"
)

func bytesValue(context ruleContext) (reflect.Value, error) {
	value := context.Value
	if !value.IsValid() {
		return reflect.Value{}, errors.New("should receive valid bytes")
	}
	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
		return reflect.Value{}, errors.New("should be a byte slice or byte array")
	}
	if value.Type().Elem().Kind() != reflect.Uint8 {
		return reflect.Value{}, errors.New("should be a byte slice or byte array")
	}
	return value, nil
}

func bytesRule(validate func(reflect.Value) error) Rule {
	return func(context ruleContext) error {
		value, err := bytesValue(context)
		if err != nil {
			return err
		}
		return validate(value)
	}
}

func bytesPredicateRule(predicate func(reflect.Value) bool, message string) Rule {
	return bytesRule(func(value reflect.Value) error {
		if !predicate(value) {
			return errors.New(message)
		}
		return nil
	})
}

func byteSlice(value reflect.Value) []byte {
	result := make([]byte, value.Len())
	for index := 0; index < value.Len(); index++ {
		result[index] = byte(value.Index(index).Uint())
	}
	return result
}
