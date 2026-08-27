package govalid

import "testing"

func TestStringLengthUsesRunes(t *testing.T) {
	if err := runRule("ação", StringLength(4)); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
	if err := runRule("ação", StringMaxLength(3)); err == nil {
		t.Fatal("expected maximum length error")
	}
}

func TestStringContentRules(t *testing.T) {
	rules := []Rule{
		StringRequired(), StringAlpha(), StringLowercase(), StringTrimmed(),
		StringStartsWith("go"), StringEndsWith("pher"), StringContains("oph"),
		StringOneOf("gopher", "codex"),
	}
	for _, rule := range rules {
		if err := runRule("gopher", rule); err != nil {
			t.Fatalf("unexpected validation error: %v", err)
		}
	}
}

func TestStringFormatRules(t *testing.T) {
	tests := []struct {
		value string
		rule  Rule
	}{
		{"user@example.com", StringEmail()},
		{"https://example.com/path", StringURL()},
		{"550e8400-e29b-41d4-a716-446655440000", StringUUID()},
	}
	for _, test := range tests {
		if err := runRule(test.value, test.rule); err != nil {
			t.Fatalf("unexpected validation error for %q: %v", test.value, err)
		}
	}
	if err := runRule("not-an-email", StringEmail()); err == nil {
		t.Fatal("expected invalid email error")
	}
}

func TestBooleanRules(t *testing.T) {
	if err := runRule(true, BoolTrue()); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
	if err := runRule(false, BoolFalse()); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
	if err := runRule(true, BoolEqual(false)); err == nil {
		t.Fatal("expected boolean equality error")
	}
	if err := runRule("true", BoolTrue()); err == nil {
		t.Fatal("expected boolean type error")
	}
}
