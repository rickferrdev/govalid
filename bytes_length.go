package govalid

import (
	"errors"
	"fmt"
	"reflect"
)

// BytesLength returns a byte-sequence validation rule for length.
func BytesLength(expect int) Rule {
	return bytesLengthRule(expect,
		func(length int) bool { return length == expect },
		fmt.Sprintf("should contain exactly %d bytes", expect))
}

// BytesMinLength returns a byte-sequence validation rule for min length.
func BytesMinLength(expect int) Rule {
	return bytesLengthRule(expect,
		func(length int) bool { return length >= expect },
		fmt.Sprintf("should contain at least %d bytes", expect))
}

// BytesMaxLength returns a byte-sequence validation rule for max length.
func BytesMaxLength(expect int) Rule {
	return bytesLengthRule(expect,
		func(length int) bool { return length <= expect },
		fmt.Sprintf("should contain at most %d bytes", expect))
}

// BytesLengthBetween returns a byte-sequence validation rule for length between.
func BytesLengthBetween(minimum, maximum int) Rule {
	return bytesRule(func(value reflect.Value) error {
		if err := validateBytesLengthRange(minimum, maximum); err != nil {
			return err
		}
		if value.Len() < minimum || value.Len() > maximum {
			return fmt.Errorf("bytes length should be between %d and %d", minimum, maximum)
		}
		return nil
	})
}

// BytesLengthNotBetween returns a byte-sequence validation rule for length not between.
func BytesLengthNotBetween(minimum, maximum int) Rule {
	return bytesRule(func(value reflect.Value) error {
		if err := validateBytesLengthRange(minimum, maximum); err != nil {
			return err
		}
		if value.Len() >= minimum && value.Len() <= maximum {
			return fmt.Errorf("bytes length should not be between %d and %d", minimum, maximum)
		}
		return nil
	})
}

func bytesLengthRule(expect int, valid func(int) bool, message string) Rule {
	return bytesRule(func(value reflect.Value) error {
		if expect < 0 {
			return errors.New("bytes length cannot be negative")
		}
		if !valid(value.Len()) {
			return errors.New(message)
		}
		return nil
	})
}

func validateBytesLengthRange(minimum, maximum int) error {
	if minimum < 0 || maximum < 0 {
		return errors.New("bytes length cannot be negative")
	}
	if minimum > maximum {
		return errors.New("minimum cannot be greater than maximum")
	}
	return nil
}
