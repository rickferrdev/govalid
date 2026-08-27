package govalid

import (
	"fmt"
	"math"
)

func FloatBits(bits int) Rule {
	return floatPredicateRule(
		func(value floatNumber) bool { return value.bits == bits },
		fmt.Sprintf("should be a %d-bit float type", bits),
	)
}

func FloatIs32() Rule { return FloatBits(32) }
func FloatIs64() Rule { return FloatBits(64) }

func FloatNaN() Rule {
	return floatPredicateRule(
		func(value floatNumber) bool { return math.IsNaN(value.value) },
		"should be NaN",
	)
}

func FloatNotNaN() Rule {
	return floatPredicateRule(
		func(value floatNumber) bool { return !math.IsNaN(value.value) },
		"should not be NaN",
	)
}

func FloatInfinite() Rule {
	return floatPredicateRule(
		func(value floatNumber) bool { return math.IsInf(value.value, 0) },
		"should be infinite",
	)
}

func FloatPositiveInfinite() Rule {
	return floatPredicateRule(
		func(value floatNumber) bool { return math.IsInf(value.value, 1) },
		"should be positive infinity",
	)
}

func FloatNegativeInfinite() Rule {
	return floatPredicateRule(
		func(value floatNumber) bool { return math.IsInf(value.value, -1) },
		"should be negative infinity",
	)
}

func FloatNotInfinite() Rule {
	return floatPredicateRule(
		func(value floatNumber) bool { return !math.IsInf(value.value, 0) },
		"should not be infinite",
	)
}

func FloatFinite() Rule {
	return floatPredicateRule(
		func(value floatNumber) bool {
			return !math.IsNaN(value.value) && !math.IsInf(value.value, 0)
		},
		"should be finite",
	)
}

func FloatNotFinite() Rule {
	return floatPredicateRule(
		func(value floatNumber) bool {
			return math.IsNaN(value.value) || math.IsInf(value.value, 0)
		},
		"should not be finite",
	)
}
