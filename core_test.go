package govalid

import (
	"errors"
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
