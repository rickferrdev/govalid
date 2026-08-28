package govalid

import (
	"fmt"
	"math"
)

// FloatBits returns a floating-point validation rule for bits.
func FloatBits(bits int) Rule {
	return floatPredicateRule(
		func(value floatNumber) bool { return value.bits == bits },
		fmt.Sprintf("should be a %d-bit float type", bits),
	)
}

// FloatIs32 returns a floating-point validation rule for is32.
func FloatIs32() Rule { return FloatBits(32) }

// FloatIs64 returns a floating-point validation rule for is64.
func FloatIs64() Rule { return FloatBits(64) }

// FloatNaN returns a floating-point validation rule for na n.
func FloatNaN() Rule {
	return floatPredicateRule(
		func(value floatNumber) bool { return math.IsNaN(value.value) },
		"should be NaN",
	)
}

// FloatNotNaN returns a floating-point validation rule for not na n.
func FloatNotNaN() Rule {
	return floatPredicateRule(
		func(value floatNumber) bool { return !math.IsNaN(value.value) },
		"should not be NaN",
	)
}

// FloatInfinite returns a floating-point validation rule for infinite.
func FloatInfinite() Rule {
	return floatPredicateRule(
		func(value floatNumber) bool { return math.IsInf(value.value, 0) },
		"should be infinite",
	)
}

// FloatPositiveInfinite returns a floating-point validation rule for positive infinite.
func FloatPositiveInfinite() Rule {
	return floatPredicateRule(
		func(value floatNumber) bool { return math.IsInf(value.value, 1) },
		"should be positive infinity",
	)
}

// FloatNegativeInfinite returns a floating-point validation rule for negative infinite.
func FloatNegativeInfinite() Rule {
	return floatPredicateRule(
		func(value floatNumber) bool { return math.IsInf(value.value, -1) },
		"should be negative infinity",
	)
}

// FloatNotInfinite returns a floating-point validation rule for not infinite.
func FloatNotInfinite() Rule {
	return floatPredicateRule(
		func(value floatNumber) bool { return !math.IsInf(value.value, 0) },
		"should not be infinite",
	)
}

// FloatFinite returns a floating-point validation rule for finite.
func FloatFinite() Rule {
	return floatPredicateRule(
		func(value floatNumber) bool {
			return !math.IsNaN(value.value) && !math.IsInf(value.value, 0)
		},
		"should be finite",
	)
}

// FloatNotFinite returns a floating-point validation rule for not finite.
func FloatNotFinite() Rule {
	return floatPredicateRule(
		func(value floatNumber) bool {
			return math.IsNaN(value.value) || math.IsInf(value.value, 0)
		},
		"should not be finite",
	)
}
