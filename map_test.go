package govalid

import "testing"

func TestMapBasicAndLengthRules(t *testing.T) {
	value := map[string]int{"one": 1, "two": 2}
	rules := []Rule{
		Map(), MapNotNil(), MapNotEmpty(), MapLength(2), MapMinLength(1),
		MapMaxLength(2), MapLengthBetween(1, 3), MapLengthNotBetween(3, 5),
	}
	for _, rule := range rules {
		if err := runRule(value, rule); err != nil {
			t.Fatalf("unexpected validation error: %v", err)
		}
	}
}

func TestMapKeyAndValueRules(t *testing.T) {
	value := map[string]string{"first": "alpha", "second": "beta"}
	rules := []Rule{
		MapHasKey("first"), MapNotHasKey("missing"),
		MapHasAllKeys("first", "second"), MapHasAnyKey("missing", "second"),
		MapHasNoneOfKeys("third"), MapAllowedKeys("first", "second", "third"),
		MapHasExactKeys("second", "first"),
		MapContainsValue("alpha"), MapNotContainsValue("gamma"),
		MapContainsAnyValue("missing", "beta"), MapContainsAllValues("alpha", "beta"),
	}
	for _, rule := range rules {
		if err := runRule(value, rule); err != nil {
			t.Fatalf("unexpected validation error: %v", err)
		}
	}
	if err := runRule(value, MapHasKey(1)); err == nil {
		t.Fatal("expected incompatible key type error")
	}
}

func TestMapNestedRulesAndComparisons(t *testing.T) {
	value := map[string]int{"one": 1, "two": 2}
	rules := []Rule{
		MapKeys(StringLowercase()), MapValues(IntPositive()), MapValueAt("two", IntEqual(2)),
		MapValueAtIfPresent("missing", IntPositive()),
		MapEqual(map[string]int{"one": 1, "two": 2}),
		MapSubsetOf(map[string]int{"one": 1, "two": 2, "three": 3}),
		MapSupersetOf(map[string]int{"one": 1}),
	}
	for _, rule := range rules {
		if err := runRule(value, rule); err != nil {
			t.Fatalf("unexpected validation error: %v", err)
		}
	}
}

func TestMapNilAndZeroValues(t *testing.T) {
	var nilMap map[string]int
	if err := runRule(nilMap, MapNil()); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
	if err := runRule(map[string]*int{"nil": nil}, MapNoNilValues()); err == nil {
		t.Fatal("expected nil value error")
	}
	if err := runRule(map[string]int{"zero": 0}, MapNoZeroValues()); err == nil {
		t.Fatal("expected zero value error")
	}
}
