package govalid

import (
	"errors"
	"fmt"
	"math"
)

// FloatGreaterThan returns a floating-point validation rule for greater than.
func FloatGreaterThan[T floating](expect T) Rule {
	expected := float64(expect)
	return floatPredicateRule(
		func(value floatNumber) bool { return value.value > expected },
		fmt.Sprintf("should be greater than %v", expected),
	)
}

// FloatGreaterThanOrEqual returns a floating-point validation rule for greater than or equal.
func FloatGreaterThanOrEqual[T floating](expect T) Rule {
	expected := float64(expect)
	return floatPredicateRule(
		func(value floatNumber) bool { return value.value >= expected },
		fmt.Sprintf("should be greater than or equal to %v", expected),
	)
}

// FloatMin returns a floating-point validation rule for min.
func FloatMin[T floating](expect T) Rule { return FloatGreaterThanOrEqual(expect) }

// FloatLessThan returns a floating-point validation rule for less than.
func FloatLessThan[T floating](expect T) Rule {
	expected := float64(expect)
	return floatPredicateRule(
		func(value floatNumber) bool { return value.value < expected },
		fmt.Sprintf("should be less than %v", expected),
	)
}

// FloatLessThanOrEqual returns a floating-point validation rule for less than or equal.
func FloatLessThanOrEqual[T floating](expect T) Rule {
	expected := float64(expect)
	return floatPredicateRule(
		func(value floatNumber) bool { return value.value <= expected },
		fmt.Sprintf("should be less than or equal to %v", expected),
	)
}

// FloatMax returns a floating-point validation rule for max.
func FloatMax[T floating](expect T) Rule { return FloatLessThanOrEqual(expect) }

// FloatBetween returns a floating-point validation rule for between.
func FloatBetween[Min floating, Max floating](minimum Min, maximum Max) Rule {
	minimumValue := float64(minimum)
	maximumValue := float64(maximum)

	return floatRule(func(value floatNumber) error {
		if math.IsNaN(minimumValue) || math.IsNaN(maximumValue) {
			return errors.New("minimum and maximum cannot be NaN")
		}
		if minimumValue > maximumValue {
			return errors.New("minimum cannot be greater than maximum")
		}
		if value.value < minimumValue || value.value > maximumValue {
			return fmt.Errorf("should be between %v and %v", minimumValue, maximumValue)
		}
		return nil
	})
}

// FloatNotBetween returns a floating-point validation rule for not between.
func FloatNotBetween[Min floating, Max floating](minimum Min, maximum Max) Rule {
	minimumValue := float64(minimum)
	maximumValue := float64(maximum)

	return floatRule(func(value floatNumber) error {
		if math.IsNaN(minimumValue) || math.IsNaN(maximumValue) {
			return errors.New("minimum and maximum cannot be NaN")
		}
		if minimumValue > maximumValue {
			return errors.New("minimum cannot be greater than maximum")
		}
		if value.value >= minimumValue && value.value <= maximumValue {
			return fmt.Errorf("should not be between %v and %v", minimumValue, maximumValue)
		}
		return nil
	})
}

// FloatEqual returns a floating-point validation rule for equal.
func FloatEqual[T floating](expect T) Rule {
	expected := float64(expect)
	return floatPredicateRule(
		func(value floatNumber) bool { return value.value == expected },
		fmt.Sprintf("should be equal to %v", expected),
	)
}

// FloatNotEqual returns a floating-point validation rule for not equal.
func FloatNotEqual[T floating](expect T) Rule {
	expected := float64(expect)
	return floatPredicateRule(
		func(value floatNumber) bool { return value.value != expected },
		fmt.Sprintf("should be different from %v", expected),
	)
}

// FloatPositive returns a floating-point validation rule for positive.
func FloatPositive() Rule { return FloatGreaterThan(0.0) }

// FloatNegative returns a floating-point validation rule for negative.
func FloatNegative() Rule { return FloatLessThan(0.0) }

// FloatNonPositive returns a floating-point validation rule for non positive.
func FloatNonPositive() Rule { return FloatLessThanOrEqual(0.0) }

// FloatNonNegative returns a floating-point validation rule for non negative.
func FloatNonNegative() Rule { return FloatGreaterThanOrEqual(0.0) }

// FloatZero returns a floating-point validation rule for zero.
func FloatZero() Rule { return FloatEqual(0.0) }

// FloatNonZero returns a floating-point validation rule for non zero.
func FloatNonZero() Rule { return FloatNotEqual(0.0) }

// FloatEqualWithin returns a floating-point validation rule for equal within.
func FloatEqualWithin[T floating](expect T, tolerance float64) Rule {
	expected := float64(expect)

	return floatRule(func(value floatNumber) error {
		if tolerance < 0 || math.IsNaN(tolerance) || math.IsInf(tolerance, 0) {
			return errors.New("tolerance should be a finite non-negative number")
		}
		if math.IsNaN(value.value) || math.IsNaN(expected) {
			return errors.New("NaN values cannot be approximately equal")
		}
		if value.value == expected {
			return nil
		}
		if math.IsInf(value.value, 0) || math.IsInf(expected, 0) ||
			math.Abs(value.value-expected) > tolerance {
			return fmt.Errorf("should be approximately %v with tolerance %v", expected, tolerance)
		}
		return nil
	})
}
