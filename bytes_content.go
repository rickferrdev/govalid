package govalid

import (
	"fmt"
)

// BytesNoZero returns a byte-sequence validation rule for no zero.
func BytesNoZero() Rule {
	return bytesSlicePredicateRule(
		func(value []byte) bool { return countByte(value, 0) == 0 },
		"bytes should not contain a zero byte")
}

// BytesHasZero returns a byte-sequence validation rule for has zero.
func BytesHasZero() Rule {
	return bytesSlicePredicateRule(
		func(value []byte) bool { return countByte(value, 0) > 0 },
		"bytes should contain a zero byte")
}

// BytesAllZero returns a byte-sequence validation rule for all zero.
func BytesAllZero() Rule {
	return bytesSlicePredicateRule(
		func(value []byte) bool { return countByte(value, 0) == len(value) },
		"all bytes should be zero")
}

// BytesNotAllZero returns a byte-sequence validation rule for not all zero.
func BytesNotAllZero() Rule {
	return bytesSlicePredicateRule(
		func(value []byte) bool { return len(value) > 0 && countByte(value, 0) != len(value) },
		"bytes should contain at least one non-zero byte")
}

// BytesUnique returns a byte-sequence validation rule for unique.
func BytesUnique() Rule {
	return bytesSliceRule(func(value []byte) error {
		seen := [256]bool{}
		for index, current := range value {
			if seen[current] {
				return fmt.Errorf("bytes should be unique; duplicate %d at index %d", current, index)
			}
			seen[current] = true
		}
		return nil
	})
}

// BytesASCII returns a byte-sequence validation rule for ascii.
func BytesASCII() Rule {
	return bytesSlicePredicateRule(func(value []byte) bool {
		for _, current := range value {
			if current > 127 {
				return false
			}
		}
		return true
	}, "bytes should contain only ASCII characters")
}

// BytesPrintableASCII returns a byte-sequence validation rule for printable ascii.
func BytesPrintableASCII() Rule {
	return bytesSlicePredicateRule(func(value []byte) bool {
		for _, current := range value {
			if current < 32 || current > 126 {
				return false
			}
		}
		return true
	}, "bytes should contain only printable ASCII characters")
}

func countByte(value []byte, expected byte) int {
	count := 0
	for _, current := range value {
		if current == expected {
			count++
		}
	}
	return count
}
