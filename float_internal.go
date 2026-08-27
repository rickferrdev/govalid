package govalid

import (
	"errors"
	"reflect"
)

func floatRule(validate func(value floatNumber) error) Rule {
	return func(context ruleContext) error {
		value, err := floatValue(context)
		if err != nil {
			return err
		}
		return validate(value)
	}
}

func floatValue(context ruleContext) (floatNumber, error) {
	return floatValueFromReflect(context.Value)
}

func floatValueFromReflect(value reflect.Value) (floatNumber, error) {
	if !value.IsValid() {
		return floatNumber{}, errors.New("should receive a valid float")
	}
	switch value.Kind() {
	case reflect.Float32, reflect.Float64:
		return floatNumber{value: value.Float(), bits: value.Type().Bits()}, nil
	default:
		return floatNumber{}, errors.New("should receive a float internally")
	}
}

func floatPredicateRule(compare func(value floatNumber) bool, message string) Rule {
	return floatRule(func(value floatNumber) error {
		if !compare(value) {
			return errors.New(message)
		}
		return nil
	})
}
