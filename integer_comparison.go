package govalid

import (
	"errors"
	"fmt"
)

func IntMin[T integer](expect T) Rule {
	expected := numberFrom(expect)
	return integerPredicateRule(
		func(value number) bool { return value.compare(expected) >= 0 },
		fmt.Sprintf("should be greater than or equal to %s", expected),
	)
}

func IntMax[T integer](expect T) Rule {
	expected := numberFrom(expect)
	return integerPredicateRule(
		func(value number) bool { return value.compare(expected) <= 0 },
		fmt.Sprintf("should be less than or equal to %s", expected),
	)
}

func IntGreaterThanOrEqual[T integer](expect T) Rule { return IntMin(expect) }
func IntLessThanOrEqual[T integer](expect T) Rule    { return IntMax(expect) }

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

// IntNoneof is kept for backward compatibility. Use IntNoneOf instead.
func IntNoneof[T integer](restricted ...T) Rule { return IntNoneOf(restricted...) }

func IntGreaterThan[T integer](expect T) Rule {
	expected := numberFrom(expect)
	return integerPredicateRule(
		func(value number) bool { return value.compare(expected) > 0 },
		fmt.Sprintf("should be strictly greater than %s", expected),
	)
}

func IntLessThan[T integer](expect T) Rule {
	expected := numberFrom(expect)
	return integerPredicateRule(
		func(value number) bool { return value.compare(expected) < 0 },
		fmt.Sprintf("should be strictly less than %s", expected),
	)
}

func IntEqual[T integer](expect T) Rule {
	expected := numberFrom(expect)
	return integerPredicateRule(
		func(value number) bool { return value.compare(expected) == 0 },
		fmt.Sprintf("should be of equal value %s", expected),
	)
}

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
