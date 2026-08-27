package govalid

import (
	"errors"
	"fmt"
	"reflect"
)

func collectionValue(context ruleContext) (reflect.Value, error) {
	value := context.Value
	if !value.IsValid() {
		return reflect.Value{}, errors.New("should receive a valid collection")
	}
	if value.Kind() != reflect.Array && value.Kind() != reflect.Slice {
		return reflect.Value{}, errors.New("should be of the slice or array type")
	}
	return value, nil
}

func collectionRule(validate func(value reflect.Value) error) Rule {
	return func(context ruleContext) error {
		value, err := collectionValue(context)
		if err != nil {
			return err
		}
		return validate(value)
	}
}

func collectionPredicateRule(compare func(value reflect.Value) bool, message string) Rule {
	return collectionRule(func(value reflect.Value) error {
		if !compare(value) {
			return errors.New(message)
		}
		return nil
	})
}

func collectionItemContext(context ruleContext, item reflect.Value, index int) ruleContext {
	return ruleContext{
		Root:  context.Root,
		Path:  fmt.Sprintf("%s[%d]", context.Path, index),
		Value: indirect(item),
	}
}
