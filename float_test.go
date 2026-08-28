package govalid

import (
	"math"
	"testing"
)

func TestFloatComparisonsAndApproximation(t *testing.T) {
	rules := []Rule{
		FloatMin(float32(1.5)),
		FloatMax(float64(2.5)),
		FloatBetween(float32(1), float64(3)),
		FloatNotBetween(float32(3), float64(5)),
		FloatEqualWithin(2.001, 0.01),
	}
	for _, rule := range rules {
		if err := runRule(float64(2), rule); err != nil {
			t.Fatalf("unexpected validation error: %v", err)
		}
	}

	if err := runRule(1.0, FloatEqualWithin(1.0, -1)); err == nil {
		t.Fatal("expected invalid tolerance error")
	}
}

func TestFloatSpecialValues(t *testing.T) {
	passing := []struct {
		value float64
		rule  Rule
	}{
		{math.NaN(), FloatNaN()},
		{math.Inf(1), FloatInfinite()},
		{math.Inf(-1), FloatInfinite()},
		{math.Inf(1), FloatPositiveInfinite()},
		{math.Inf(-1), FloatNegativeInfinite()},
		{1.5, FloatFinite()},
		{math.NaN(), FloatNotFinite()},
	}
	for _, test := range passing {
		if err := runRule(test.value, test.rule); err != nil {
			t.Fatalf("unexpected validation error: %v", err)
		}
	}

	if err := runRule(math.NaN(), FloatFinite()); err == nil {
		t.Fatal("expected NaN to be non-finite")
	}
	if err := runRule(math.Inf(-1), FloatNotInfinite()); err == nil {
		t.Fatal("expected negative infinity to fail FloatNotInfinite")
	}
}

func TestFloatBits(t *testing.T) {
	if err := runRule(float32(1), FloatIs32()); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
	if err := runRule(float64(1), FloatIs64()); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}
