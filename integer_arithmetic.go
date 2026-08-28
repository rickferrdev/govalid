package govalid

import (
	"errors"
	"fmt"
)

// IntEven returns an integer validation rule for even.
func IntEven() Rule {
	return integerRule(func(value number) error {
		if value.magnitude()%2 != 0 {
			return errors.New("should be an even number")
		}
		return nil
	})
}

// IntOdd returns an integer validation rule for odd.
func IntOdd() Rule {
	return integerRule(func(value number) error {
		if value.magnitude()%2 == 0 {
			return errors.New("should be an odd number")
		}
		return nil
	})
}

// IntMultipleOf returns an integer validation rule for multiple of.
func IntMultipleOf[T integer](n T) Rule {
	return divisibleIntegerRule(numberFrom(n), "should be a multiple of")
}

// IntDivisibleBy returns an integer validation rule for divisible by.
func IntDivisibleBy[T integer](n T) Rule {
	return divisibleIntegerRule(numberFrom(n), "should be divisible by")
}

// IntPrime returns an integer validation rule for prime.
func IntPrime() Rule {
	return integerRule(func(value number) error {
		if value.isNegative() || !isPrime(value.magnitude()) {
			return errors.New("should be a prime number")
		}
		return nil
	})
}

// IntCompound returns an integer validation rule for compound.
func IntCompound() Rule {
	return integerRule(func(value number) error {
		if value.isNegative() || value.magnitude() < 4 || isPrime(value.magnitude()) {
			return errors.New("should be a compound number")
		}
		return nil
	})
}

// IntPowerOfTwo returns an integer validation rule for power of two.
func IntPowerOfTwo() Rule {
	return integerRule(func(value number) error {
		magnitude := value.magnitude()
		if value.isNegative() || magnitude == 0 || magnitude&(magnitude-1) != 0 {
			return errors.New("should be a power of 2")
		}
		return nil
	})
}

// IntPerfectSquare returns an integer validation rule for perfect square.
func IntPerfectSquare() Rule {
	return integerRule(func(value number) error {
		if value.isNegative() || !isPerfectSquare(value.magnitude()) {
			return errors.New("should be a perfect square")
		}
		return nil
	})
}

func divisibleIntegerRule(divisor number, message string) Rule {
	return integerRule(func(value number) error {
		magnitude := divisor.magnitude()
		if magnitude == 0 {
			return errors.New("divisor cannot be zero")
		}
		if value.magnitude()%magnitude != 0 {
			return fmt.Errorf("%s %s %s", value, message, divisor)
		}
		return nil
	})
}

func isPrime(value uint64) bool {
	if value < 2 {
		return false
	}
	if value%2 == 0 {
		return value == 2
	}
	for divisor := uint64(3); divisor <= value/divisor; divisor += 2 {
		if value%divisor == 0 {
			return false
		}
	}
	return true
}

func isPerfectSquare(value uint64) bool {
	if value < 2 {
		return true
	}
	low, high := uint64(1), uint64(1)<<32
	for low <= high {
		root := low + (high-low)/2
		quotient := value / root
		switch {
		case root == quotient && value%root == 0:
			return true
		case root > quotient:
			high = root - 1
		default:
			low = root + 1
		}
	}
	return false
}
