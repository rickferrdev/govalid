package govalid

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

type validationFixture struct {
	Name   string
	Age    int
	Active bool
	Nested struct{ Code string }
}

func TestValidateAcceptsValidStruct(t *testing.T) {
	data := validationFixture{Name: "Alice", Age: 20, Active: true}
	data.Nested.Code = "ABC"

	err := New().Validate(data,
		Field("Name", StringRequired(), StringMinLength(3)),
		Field("Age", IntBetween(18, 120)),
		Field("Active", BoolTrue()),
		Field("Nested.Code", StringLength(3)),
	)
	if err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestValidateCollectsIssues(t *testing.T) {
	data := validationFixture{Name: "", Age: 10}
	err := New().Validate(data,
		Field("Name", StringRequired()),
		Field("Age", IntMin(18)),
		Field("Active", BoolTrue()),
	)

	var validationErr *FieldIssueError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected FieldIssueError, received %T", err)
	}
	if len(validationErr.RulesIssues) != 3 {
		t.Fatalf("expected 3 issues, received %d", len(validationErr.RulesIssues))
	}
}

func TestValidateStopsOnFirstIssue(t *testing.T) {
	data := validationFixture{}
	err := New(WithStopOnFirstError()).Validate(data,
		Field("Name", StringRequired()),
		Field("Age", IntPositive()),
	)

	var validationErr *FieldIssueError
	if !errors.As(err, &validationErr) || len(validationErr.RulesIssues) != 1 {
		t.Fatalf("expected exactly one issue, received %v", err)
	}
}

func TestValidateRejectsInvalidInputAndPath(t *testing.T) {
	if err := New().Validate(10); err == nil {
		t.Fatal("expected non-struct input error")
	}
	if err := New().Validate(validationFixture{}, Field("Missing", StringRequired())); err == nil {
		t.Fatal("expected missing field error")
	}
}

func TestValidatePanicsWithRuleContext(t *testing.T) {
	defer func() {
		recovered := recover()
		panicErr, ok := recovered.(*Issue)
		if !ok {
			t.Fatalf("expected *RulePanicError, received %T", recovered)
		}

		if panicErr.Path != "Age" || panicErr.Value != 10 {
			t.Fatalf("unexpected panic context: %#v", panicErr)
		}
		if !strings.Contains(panicErr.Rule, "IntMin") {
			t.Fatalf("expected rule name to contain IntMin, received %q", panicErr.Rule)
		}

		message := fmt.Sprint(panicErr)
		for _, expected := range []string{"rule:", "path: Age", "value: 10 (int)", "message: should be greater than or equal to 18"} {
			if !strings.Contains(message, expected) {
				t.Fatalf("panic message %q does not contain %q", message, expected)
			}
		}
	}()

	New(WithPanicOnFirstError()).Validate(
		validationFixture{Age: 10},
		Field("Age", IntMin(18)),
	)
}

func TestIssueHandlerReceivesSilencedIssues(t *testing.T) {
	var handled []Issue
	validator := New(
		WithIssueHandler(func(issue Issue) { handled = append(handled, issue) }),
		WithSilenceErrors(),
		WithStopOnFirstError(),
	)

	err := validator.Validate(validationFixture{},
		Field("Name", StringRequired()),
		Field("Age", IntPositive()),
	)
	if err != nil {
		t.Fatalf("silenced validation returned an error: %v", err)
	}
	if len(handled) != 1 {
		t.Fatalf("expected handler to receive 1 issue, received %d", len(handled))
	}
	if handled[0].Path != "Name" {
		t.Fatalf("unexpected handled issues: %#v", handled)
	}
}
