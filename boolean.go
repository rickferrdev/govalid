package govalid

import (
	"errors"
	"fmt"
	"reflect"
)

// BoolTrue returns a boolean validation rule for true.
func BoolTrue() Rule {
	return boolRule(func(b bool) error {
		if !b {
			return errors.New("should be true")
		}

		return nil
	})
}

// BoolFalse returns a boolean validation rule for false.
func BoolFalse() Rule {
	return boolRule(func(b bool) error {
		if b {
			return errors.New("should be false")
		}

		return nil
	})
}

// BoolEqual returns a boolean validation rule for equal.
func BoolEqual(expect bool) Rule {
	return boolRule(func(b bool) error {
		if b != expect {
			return fmt.Errorf(" should be %v", expect)
		}

		return nil
	})
}

func boolRule(validate func(bool) error) Rule {
	return func(context ruleContext) error {
		value, err := boolValue(context)
		if err != nil {
			return err
		}
		return validate(value)
	}
}

func boolValue(context ruleContext) (bool, error) {
	if !context.Value.IsValid() || context.Value.Kind() != reflect.Bool {
		return false, errors.New("should be of type boolean")
	}
	return context.Value.Bool(), nil
}
