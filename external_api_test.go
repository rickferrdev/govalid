package govalid_test

import (
	"testing"

	"github.com/rickferrdev/govalid"
)

func TestGenericRulesRemainCallableByConsumers(t *testing.T) {
	type customInt uint16
	type customFloat float32

	_ = govalid.IntMin(customInt(10))
	_ = govalid.IntBetween(int8(-1), uint64(100))
	_ = govalid.FloatMin(customFloat(1.5))
	_ = govalid.FloatBetween(float32(1), float64(2))
}
