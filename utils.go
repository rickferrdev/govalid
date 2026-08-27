package govalid

import (
	"fmt"
	"reflect"
	"strings"
)

func lookup(root reflect.Value, path string) (reflect.Value, error) {
	current := root

	for part := range strings.SplitSeq(path, ".") {
		current = indirect(current)

		if current.Kind() != reflect.Struct {
			return reflect.Value{}, fmt.Errorf(
				"%q cannot be acessed on %s",
				part,
				current.Kind(),
			)
		}

		current = current.FieldByName(part)
		if !current.IsValid() {
			return reflect.Value{}, fmt.Errorf("field %q was not found", path)
		}
	}

	return indirect(current), nil
}

func indirect(value reflect.Value) reflect.Value {
	for value.IsValid() && (value.Kind() == reflect.Pointer || value.Kind() == reflect.Interface) {
		if value.IsNil() {
			return value
		}

		value = value.Elem()
	}

	return value
}

func applyRules(context ruleContext, rules []Rule) error {
	for _, rule := range rules {
		if rule == nil {
			continue
		}
		if err := rule(context); err != nil {
			return err
		}
	}
	return nil
}
