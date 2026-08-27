package govalid

import (
	"fmt"
)

func BytesNoZero() Rule {
	return bytesSlicePredicateRule(
		func(value []byte) bool { return countByte(value, 0) == 0 },
		"bytes should not contain a zero byte")
}

func BytesHasZero() Rule {
	return bytesSlicePredicateRule(
		func(value []byte) bool { return countByte(value, 0) > 0 },
		"bytes should contain a zero byte")
}

func BytesAllZero() Rule {
	return bytesSlicePredicateRule(
		func(value []byte) bool { return countByte(value, 0) == len(value) },
		"all bytes should be zero")
}

func BytesNotAllZero() Rule {
	return bytesSlicePredicateRule(
		func(value []byte) bool { return len(value) > 0 && countByte(value, 0) != len(value) },
		"bytes should contain at least one non-zero byte")
}

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
