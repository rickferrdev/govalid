package govalid

import (
	"math"
	"testing"
)

func TestIntegerComparisonsSupportSignedAndUnsigned(t *testing.T) {
	tests := []struct {
		value any
		rule  Rule
	}{
		{int8(10), IntMin(uint16(10))},
		{uint32(10), IntMax(int64(10))},
		{uint64(100), IntBetween(int8(-1), uint64(100))},
		{uint64(math.MaxUint64), IntEqual(uint64(math.MaxUint64))},
		{int64(-1), IntLessThan(uint64(0))},
	}
	for _, test := range tests {
		if err := runRule(test.value, test.rule); err != nil {
			t.Fatalf("unexpected validation error for %v: %v", test.value, err)
		}
	}

	if err := runRule(uint64(math.MaxUint64), IntMax(int64(math.MaxInt64))); err == nil {
		t.Fatal("expected MaxUint64 to exceed MaxInt64")
	}
}

func TestIntegerArithmeticRules(t *testing.T) {
	tests := []struct {
		value any
		rule  Rule
	}{
		{int64(-3), IntOdd()},
		{uint64(12), IntEven()},
		{uint64(12), IntDivisibleBy(uint8(3))},
		{int64(-12), IntMultipleOf(int8(-3))},
		{uint64(13), IntPrime()},
		{uint64(12), IntCompound()},
		{uint64(1) << 63, IntPowerOfTwo()},
		{uint64(1) << 62, IntPerfectSquare()},
	}
	for _, test := range tests {
		if err := runRule(test.value, test.rule); err != nil {
			t.Fatalf("unexpected validation error for %v: %v", test.value, err)
		}
	}

	if err := runRule(12, IntDivisibleBy(0)); err == nil {
		t.Fatal("expected zero divisor error")
	}
}

func TestIntegerSetRules(t *testing.T) {
	if err := runRule(2, IntOneOf(1, 2, 3)); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
	if err := runRule(4, IntNoneOf(1, 2, 3)); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
	if err := runRule(2, IntNotBetween(1, 3)); err == nil {
		t.Fatal("expected value inside forbidden interval to fail")
	}
}
