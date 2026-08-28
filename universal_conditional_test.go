package govalid

import "testing"

func TestUniversalRules(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		rule      Rule
		shouldErr bool
	}{
		{name: "required string", value: "value", rule: Required()},
		{name: "required empty string", value: "", rule: Required(), shouldErr: true},
		{name: "required zero integer", value: 0, rule: Required()},
		{name: "required false", value: false, rule: Required()},
		{name: "nil slice", value: []string(nil), rule: Nil()},
		{name: "non-nil slice", value: []string{}, rule: NotNil()},
		{name: "zero integer", value: 0, rule: Zero()},
		{name: "non-zero integer", value: 1, rule: Zero(), shouldErr: true},
		{name: "not-zero string", value: "value", rule: NotZero()},
		{name: "zero string", value: "", rule: NotZero(), shouldErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := runRule(test.value, test.rule)
			if (err != nil) != test.shouldErr {
				t.Fatalf("unexpected error state: %v", err)
			}
		})
	}
}

func TestConditionalRules(t *testing.T) {
	if err := runRule("", When(false, StringRequired())); err != nil {
		t.Fatalf("disabled When rejected value: %v", err)
	}
	if err := runRule("", When(true, StringRequired())); err == nil {
		t.Fatal("enabled When should apply nested rules")
	}
	if err := runRule("", Unless(true, StringRequired())); err != nil {
		t.Fatalf("disabled Unless rejected value: %v", err)
	}
	if err := runRule("", Unless(false, StringRequired())); err == nil {
		t.Fatal("enabled Unless should apply nested rules")
	}
	if err := runRule("", Optional(StringRequired())); err != nil {
		t.Fatalf("Optional should skip an empty value: %v", err)
	}
	if err := runRule(0, Optional(IntPositive())); err == nil {
		t.Fatal("Optional should validate present scalar zero values")
	}
}

func TestContextConditionsCanInspectRoot(t *testing.T) {
	type fixture struct {
		Enabled bool
		Name    string
	}

	requiredWhenEnabled := WhenContext(func(context RuleContext) bool {
		return context.RootAny().(fixture).Enabled
	}, StringRequired())

	validator := New()
	if err := validator.Validate(fixture{}, Field("Name", requiredWhenEnabled)); err != nil {
		t.Fatalf("disabled contextual rule rejected value: %v", err)
	}
	if err := validator.Validate(fixture{Enabled: true}, Field("Name", requiredWhenEnabled)); err == nil {
		t.Fatal("enabled contextual rule should reject an empty name")
	}
	if err := runRule("value", WhenContext(nil, StringRequired())); err == nil {
		t.Fatal("nil contextual condition should fail")
	}
	if err := runRule("", UnlessContext(nil, StringRequired())); err == nil {
		t.Fatal("nil contextual condition should fail")
	}
}
