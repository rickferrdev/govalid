package govalid

import (
	"strings"
	"testing"
)

func TestCollectionBasicAndLengthRules(t *testing.T) {
	value := []int{1, 2, 3}
	rules := []Rule{
		Collection(), CollectionNotNil(), CollectionNotEmpty(), CollectionLength(3),
		CollectionMinLength(2), CollectionMaxLength(4),
		CollectionLengthBetween(2, 4), CollectionLengthNotBetween(4, 6),
	}
	for _, rule := range rules {
		if err := runRule(value, rule); err != nil {
			t.Fatalf("unexpected validation error: %v", err)
		}
	}
	var nilSlice []int
	if err := runRule(nilSlice, CollectionNil()); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
	if err := runRule([0]int{}, CollectionNotNil()); err != nil {
		t.Fatalf("arrays should be non-nil collections: %v", err)
	}
}

func TestCollectionContentRules(t *testing.T) {
	value := []string{"alpha", "beta", "gamma"}
	rules := []Rule{
		CollectionContains("alpha"), CollectionNotContains("delta"),
		CollectionContainsAny("missing", "beta"), CollectionContainsAll("alpha", "gamma"),
		CollectionUnique(), CollectionNoNilItems(), CollectionNoZeroItems(),
		CollectionHasNonZeroItem(),
	}
	for _, rule := range rules {
		if err := runRule(value, rule); err != nil {
			t.Fatalf("unexpected validation error: %v", err)
		}
	}
	if err := runRule([]int{1, 1}, CollectionUnique()); err == nil {
		t.Fatal("expected duplicate item error")
	}
}

func TestCollectionNestedRules(t *testing.T) {
	value := []int{2, 4, 6}
	rules := []Rule{
		CollectionEach(IntPositive(), IntEven()), CollectionAny(IntEqual(4)),
		CollectionNone(IntNegative()), CollectionItemAt(1, IntEqual(4)),
		CollectionItemAtIfPresent(10, IntPositive()),
	}
	for _, rule := range rules {
		if err := runRule(value, rule); err != nil {
			t.Fatalf("unexpected validation error: %v", err)
		}
	}

	err := runRule([]string{"valid", ""}, CollectionEach(StringRequired()))
	if err == nil || !strings.Contains(err.Error(), "Value[1]") {
		t.Fatalf("expected indexed error, received %v", err)
	}
}

func TestCollectionComparisonRules(t *testing.T) {
	if err := runRule([]int{1, 2}, CollectionEqual([]int{1, 2})); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
	if err := runRule([]int{1, 2}, CollectionNotEqual([]int{2, 1})); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}
