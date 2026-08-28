package govalid

import (
	"errors"
	"fmt"
)

// IntMin returns an integer validation rule for min.
func IntMin[T integer](expect T) Rule {
	expected := numberFrom(expect)
	return integerPredicateRule(
		func(value number) bool { return value.compare(expected) >= 0 },
		fmt.Sprintf("should be greater than or equal to %s", expected),
	)
}

// IntMax returns an integer validation rule for max.
func IntMax[T integer](expect T) Rule {
	expected := numberFrom(expect)
	return integerPredicateRule(
		func(value number) bool { return value.compare(expected) <= 0 },
		fmt.Sprintf("should be less than or equal to %s", expected),
	)
}

// IntGreaterThanOrEqual returns an integer validation rule for greater than or equal.
func IntGreaterThanOrEqual[T integer](expect T) Rule { return IntMin(expect) }

// IntLessThanOrEqual returns an integer validation rule for less than or equal.
func IntLessThanOrEqual[T integer](expect T) Rule { return IntMax(expect) }

// IntBetween returns an integer validation rule for between.
func IntBetween[Min integer, Max integer](minimum Min, maximum Max) Rule {
	minimumNumber := numberFrom(minimum)
	maximumNumber := numberFrom(maximum)

	return integerRule(func(value number) error {
		if minimumNumber.compare(maximumNumber) > 0 {
			return errors.New("minimum cannot be greater than maximum")
		}
		if value.compare(minimumNumber) < 0 || value.compare(maximumNumber) > 0 {
			return fmt.Errorf("should be between %s and %s", minimumNumber, maximumNumber)
		}
		return nil
	})
}

// IntNotBetween returns an integer validation rule for not between.
func IntNotBetween[Min integer, Max integer](minimum Min, maximum Max) Rule {
	minimumNumber := numberFrom(minimum)
	maximumNumber := numberFrom(maximum)

	return integerRule(func(value number) error {
		if minimumNumber.compare(maximumNumber) > 0 {
			return errors.New("minimum cannot be greater than maximum")
		}
		if value.compare(minimumNumber) >= 0 && value.compare(maximumNumber) <= 0 {
			return fmt.Errorf("should not be between %s and %s", minimumNumber, maximumNumber)
		}
		return nil
	})
}

// IntOneOf returns an integer validation rule for one of.
func IntOneOf[T integer](allowed ...T) Rule {
	allowedNumbers := make([]number, len(allowed))
	for index, value := range allowed {
		allowedNumbers[index] = numberFrom(value)
	}

	return integerRule(func(value number) error {
		for _, allowed := range allowedNumbers {
			if value.compare(allowed) == 0 {
				return nil
			}
		}
		return fmt.Errorf("should be one of %v", allowedNumbers)
	})
}

// IntNoneOf returns an integer validation rule for none of.
func IntNoneOf[T integer](restricted ...T) Rule {
	restrictedNumbers := make([]number, len(restricted))
	for index, value := range restricted {
		restrictedNumbers[index] = numberFrom(value)
	}

	return integerRule(func(value number) error {
		for _, restricted := range restrictedNumbers {
			if value.compare(restricted) == 0 {
				return fmt.Errorf("should not be %s", restricted)
			}
		}
		return nil
	})
}

// IntGreaterThan returns an integer validation rule for greater than.
func IntGreaterThan[T integer](expect T) Rule {
	expected := numberFrom(expect)
	return integerPredicateRule(
		func(value number) bool { return value.compare(expected) > 0 },
		fmt.Sprintf("should be strictly greater than %s", expected),
	)
}

// IntLessThan returns an integer validation rule for less than.
func IntLessThan[T integer](expect T) Rule {
	expected := numberFrom(expect)
	return integerPredicateRule(
		func(value number) bool { return value.compare(expected) < 0 },
		fmt.Sprintf("should be strictly less than %s", expected),
	)
}

// IntEqual returns an integer validation rule for equal.
func IntEqual[T integer](expect T) Rule {
	expected := numberFrom(expect)
	return integerPredicateRule(
		func(value number) bool { return value.compare(expected) == 0 },
		fmt.Sprintf("should be of equal value %s", expected),
	)
}

// IntNotEqual returns an integer validation rule for not equal.
func IntNotEqual[T integer](expect T) Rule {
	expected := numberFrom(expect)
	return integerPredicateRule(
		func(value number) bool { return value.compare(expected) != 0 },
		"should be different values",
	)
}

func integerPredicateRule(valid func(number) bool, message string) Rule {
	return integerRule(func(value number) error {
		if !valid(value) {
			return errors.New(message)
		}
		return nil
	})
}
