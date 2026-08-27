package govalid

import (
	"errors"
	"reflect"
)

func stringValue(context ruleContext) (string, error) {
	if context.Value.Kind() != reflect.String {
		return "", errors.New("should be of type string")
	}
	return context.Value.String(), nil
}

func stringRuneRule(message string, valid func(rune) bool) Rule {
	return stringRule(func(value string) error {
		for _, char := range value {
			if !valid(char) {
				return errors.New(message)
			}
		}
		return nil
	})
}

func stringPredicateRule(predicate func(string) bool, message string) Rule {
	return stringRule(func(value string) error {
		if !predicate(value) {
			return errors.New(message)
		}
		return nil
	})
}

func stringRule(validate func(string) error) Rule {
	return func(context ruleContext) error {
		value, err := stringValue(context)
		if err != nil {
			return err
		}
		return validate(value)
	}
}
