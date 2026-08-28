package govalid

import (
	"bytes"
	"crypto/subtle"
	"errors"
	"fmt"
	"reflect"
)

// BytesEqual returns a byte-sequence validation rule for equal.
func BytesEqual(expected []byte) Rule {
	expected = bytes.Clone(expected)
	return bytesSliceRule(func(value []byte) error {
		if !bytes.Equal(value, expected) {
			return errors.New("bytes should equal the expected value")
		}
		return nil
	})
}

// BytesNotEqual returns a byte-sequence validation rule for not equal.
func BytesNotEqual(unexpected []byte) Rule {
	unexpected = bytes.Clone(unexpected)
	return bytesSliceRule(func(value []byte) error {
		if bytes.Equal(value, unexpected) {
			return errors.New("bytes should differ from the unexpected value")
		}
		return nil
	})
}

// BytesConstantTimeEqual returns a byte-sequence validation rule for constant time equal.
func BytesConstantTimeEqual(expected []byte) Rule {
	expected = bytes.Clone(expected)
	return bytesSliceRule(func(value []byte) error {
		if subtle.ConstantTimeCompare(value, expected) != 1 {
			return errors.New("bytes should equal the expected value")
		}
		return nil
	})
}

// BytesOneOf returns a byte-sequence validation rule for one of.
func BytesOneOf(expected ...[]byte) Rule {
	expected = cloneByteSlices(expected)
	return bytesSliceRule(func(value []byte) error {
		for _, candidate := range expected {
			if bytes.Equal(value, candidate) {
				return nil
			}
		}
		return errors.New("bytes should equal one of the expected values")
	})
}

// BytesNoneOf returns a byte-sequence validation rule for none of.
func BytesNoneOf(unexpected ...[]byte) Rule {
	unexpected = cloneByteSlices(unexpected)
	return bytesSliceRule(func(value []byte) error {
		for _, candidate := range unexpected {
			if bytes.Equal(value, candidate) {
				return errors.New("bytes should differ from every restricted value")
			}
		}
		return nil
	})
}

// BytesContains returns a byte-sequence validation rule for contains.
func BytesContains(expected []byte) Rule {
	expected = bytes.Clone(expected)
	return bytesSlicePredicateRule(
		func(value []byte) bool { return bytes.Contains(value, expected) },
		fmt.Sprintf("bytes should contain %v", expected))
}

// BytesNotContains returns a byte-sequence validation rule for not contains.
func BytesNotContains(unexpected []byte) Rule {
	unexpected = bytes.Clone(unexpected)
	return bytesSlicePredicateRule(
		func(value []byte) bool { return !bytes.Contains(value, unexpected) },
		fmt.Sprintf("bytes should not contain %v", unexpected))
}

// BytesContainsAny returns a byte-sequence validation rule for contains any.
func BytesContainsAny(expected ...[]byte) Rule {
	expected = cloneByteSlices(expected)
	return bytesSliceRule(func(value []byte) error {
		for _, candidate := range expected {
			if bytes.Contains(value, candidate) {
				return nil
			}
		}
		return errors.New("bytes should contain at least one expected sequence")
	})
}

// BytesContainsAll returns a byte-sequence validation rule for contains all.
func BytesContainsAll(expected ...[]byte) Rule {
	expected = cloneByteSlices(expected)
	return bytesSliceRule(func(value []byte) error {
		for _, candidate := range expected {
			if !bytes.Contains(value, candidate) {
				return fmt.Errorf("bytes should contain %v", candidate)
			}
		}
		return nil
	})
}

// BytesHasPrefix returns a byte-sequence validation rule for has prefix.
func BytesHasPrefix(prefix []byte) Rule {
	prefix = bytes.Clone(prefix)
	return bytesSlicePredicateRule(
		func(value []byte) bool { return bytes.HasPrefix(value, prefix) },
		fmt.Sprintf("bytes should have prefix %v", prefix))
}

// BytesNotHasPrefix returns a byte-sequence validation rule for not has prefix.
func BytesNotHasPrefix(prefix []byte) Rule {
	prefix = bytes.Clone(prefix)
	return bytesSlicePredicateRule(
		func(value []byte) bool { return !bytes.HasPrefix(value, prefix) },
		fmt.Sprintf("bytes should not have prefix %v", prefix))
}

// BytesHasSuffix returns a byte-sequence validation rule for has suffix.
func BytesHasSuffix(suffix []byte) Rule {
	suffix = bytes.Clone(suffix)
	return bytesSlicePredicateRule(
		func(value []byte) bool { return bytes.HasSuffix(value, suffix) },
		fmt.Sprintf("bytes should have suffix %v", suffix))
}

// BytesNotHasSuffix returns a byte-sequence validation rule for not has suffix.
func BytesNotHasSuffix(suffix []byte) Rule {
	suffix = bytes.Clone(suffix)
	return bytesSlicePredicateRule(
		func(value []byte) bool { return !bytes.HasSuffix(value, suffix) },
		fmt.Sprintf("bytes should not have suffix %v", suffix))
}

// BytesGreaterThan returns a byte-sequence validation rule for greater than.
func BytesGreaterThan(expected []byte) Rule {
	return bytesOrderedRule(expected, func(order int) bool { return order > 0 }, "greater than")
}

// BytesGreaterThanOrEqual returns a byte-sequence validation rule for greater than or equal.
func BytesGreaterThanOrEqual(expected []byte) Rule {
	return bytesOrderedRule(expected, func(order int) bool { return order >= 0 }, "greater than or equal to")
}

// BytesLessThan returns a byte-sequence validation rule for less than.
func BytesLessThan(expected []byte) Rule {
	return bytesOrderedRule(expected, func(order int) bool { return order < 0 }, "less than")
}

// BytesLessThanOrEqual returns a byte-sequence validation rule for less than or equal.
func BytesLessThanOrEqual(expected []byte) Rule {
	return bytesOrderedRule(expected, func(order int) bool { return order <= 0 }, "less than or equal to")
}

// BytesBetween returns a byte-sequence validation rule for between.
func BytesBetween(minimum, maximum []byte) Rule {
	minimum = bytes.Clone(minimum)
	maximum = bytes.Clone(maximum)
	return bytesSliceRule(func(value []byte) error {
		if bytes.Compare(minimum, maximum) > 0 {
			return errors.New("minimum bytes cannot be greater than maximum bytes")
		}
		if bytes.Compare(value, minimum) < 0 || bytes.Compare(value, maximum) > 0 {
			return errors.New("bytes should be inside the expected lexicographical range")
		}
		return nil
	})
}

func bytesOrderedRule(expected []byte, valid func(int) bool, relation string) Rule {
	expected = bytes.Clone(expected)
	return bytesSliceRule(func(value []byte) error {
		if !valid(bytes.Compare(value, expected)) {
			return fmt.Errorf("bytes should be %s the expected value", relation)
		}
		return nil
	})
}

func cloneByteSlices(values [][]byte) [][]byte {
	cloned := make([][]byte, len(values))
	for index, value := range values {
		cloned[index] = bytes.Clone(value)
	}
	return cloned
}

func bytesSliceRule(validate func([]byte) error) Rule {
	return bytesRule(func(value reflect.Value) error { return validate(byteSlice(value)) })
}

func bytesSlicePredicateRule(predicate func([]byte) bool, message string) Rule {
	return bytesSliceRule(func(value []byte) error {
		if !predicate(value) {
			return errors.New(message)
		}
		return nil
	})
}
