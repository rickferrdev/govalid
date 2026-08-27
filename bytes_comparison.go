package govalid

import (
	"bytes"
	"crypto/subtle"
	"errors"
	"fmt"
	"reflect"
)

func BytesEqual(expected []byte) Rule {
	expected = bytes.Clone(expected)
	return bytesSliceRule(func(value []byte) error {
		if !bytes.Equal(value, expected) {
			return errors.New("bytes should equal the expected value")
		}
		return nil
	})
}

func BytesNotEqual(unexpected []byte) Rule {
	unexpected = bytes.Clone(unexpected)
	return bytesSliceRule(func(value []byte) error {
		if bytes.Equal(value, unexpected) {
			return errors.New("bytes should differ from the unexpected value")
		}
		return nil
	})
}

func BytesConstantTimeEqual(expected []byte) Rule {
	expected = bytes.Clone(expected)
	return bytesSliceRule(func(value []byte) error {
		if subtle.ConstantTimeCompare(value, expected) != 1 {
			return errors.New("bytes should equal the expected value")
		}
		return nil
	})
}

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

func BytesContains(expected []byte) Rule {
	expected = bytes.Clone(expected)
	return bytesSlicePredicateRule(
		func(value []byte) bool { return bytes.Contains(value, expected) },
		fmt.Sprintf("bytes should contain %v", expected))
}

func BytesNotContains(unexpected []byte) Rule {
	unexpected = bytes.Clone(unexpected)
	return bytesSlicePredicateRule(
		func(value []byte) bool { return !bytes.Contains(value, unexpected) },
		fmt.Sprintf("bytes should not contain %v", unexpected))
}

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

func BytesHasPrefix(prefix []byte) Rule {
	prefix = bytes.Clone(prefix)
	return bytesSlicePredicateRule(
		func(value []byte) bool { return bytes.HasPrefix(value, prefix) },
		fmt.Sprintf("bytes should have prefix %v", prefix))
}

func BytesNotHasPrefix(prefix []byte) Rule {
	prefix = bytes.Clone(prefix)
	return bytesSlicePredicateRule(
		func(value []byte) bool { return !bytes.HasPrefix(value, prefix) },
		fmt.Sprintf("bytes should not have prefix %v", prefix))
}

func BytesHasSuffix(suffix []byte) Rule {
	suffix = bytes.Clone(suffix)
	return bytesSlicePredicateRule(
		func(value []byte) bool { return bytes.HasSuffix(value, suffix) },
		fmt.Sprintf("bytes should have suffix %v", suffix))
}

func BytesNotHasSuffix(suffix []byte) Rule {
	suffix = bytes.Clone(suffix)
	return bytesSlicePredicateRule(
		func(value []byte) bool { return !bytes.HasSuffix(value, suffix) },
		fmt.Sprintf("bytes should not have suffix %v", suffix))
}

func BytesGreaterThan(expected []byte) Rule {
	return bytesOrderedRule(expected, func(order int) bool { return order > 0 }, "greater than")
}

func BytesGreaterThanOrEqual(expected []byte) Rule {
	return bytesOrderedRule(expected, func(order int) bool { return order >= 0 }, "greater than or equal to")
}

func BytesLessThan(expected []byte) Rule {
	return bytesOrderedRule(expected, func(order int) bool { return order < 0 }, "less than")
}

func BytesLessThanOrEqual(expected []byte) Rule {
	return bytesOrderedRule(expected, func(order int) bool { return order <= 0 }, "less than or equal to")
}

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
