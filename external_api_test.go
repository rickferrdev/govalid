package govalid_test

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
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

func TestConsumersCanDefineRules(t *testing.T) {
	type payload struct {
		Name string
	}

	customRule := govalid.Rule(func(context govalid.RuleContext) error {
		if context.Path != "Name" {
			return fmt.Errorf("unexpected path %q", context.Path)
		}
		if context.ValueAny() != "accepted" {
			return errors.New("value was not accepted")
		}
		if _, ok := context.RootAny().(payload); !ok {
			return fmt.Errorf("unexpected root type %T", context.RootAny())
		}
		return nil
	})

	validator := govalid.New()
	if err := validator.Validate(payload{Name: "accepted"}, govalid.Field("Name", customRule)); err != nil {
		t.Fatalf("custom rule rejected valid input: %v", err)
	}

	err := validator.Validate(payload{Name: "rejected"}, govalid.Field("Name", customRule))
	var issue *govalid.Issue
	if !errors.As(err, &issue) {
		t.Fatalf("expected errors.As to extract *govalid.Issue from %T", err)
	}
	if issue.Path != "Name" || issue.Value != "rejected" || issue.Message != "value was not accepted" {
		t.Fatalf("unexpected extracted issue: %#v", issue)
	}
}

func TestValidatorCanBeReusedConcurrently(t *testing.T) {
	type payload struct {
		Age int
	}

	const workers = 100
	var handled atomic.Int64
	var sourceCalls atomic.Int64
	validator := govalid.New(
		govalid.WithIssueHandler(func(govalid.Issue) {
			handled.Add(1)
		}),
		govalid.WithFieldSources(govalid.FieldSourceFunc(func(reflect.Type) ([]govalid.FieldSpec, error) {
			sourceCalls.Add(1)
			return nil, nil
		})),
	)
	field := govalid.Field("Age", govalid.IntMin(18))

	var waitGroup sync.WaitGroup
	errorsChannel := make(chan error, workers)
	start := make(chan struct{})

	for index := 0; index < workers; index++ {
		waitGroup.Add(1)
		go func(age int) {
			defer waitGroup.Done()
			<-start

			err := validator.Validate(payload{Age: age}, field)
			if age >= 18 && err != nil {
				errorsChannel <- fmt.Errorf("valid age %d failed: %w", age, err)
				return
			}
			if age < 18 && err == nil {
				errorsChannel <- fmt.Errorf("invalid age %d was accepted", age)
			}
		}(index)
	}

	close(start)
	waitGroup.Wait()
	close(errorsChannel)

	for err := range errorsChannel {
		t.Error(err)
	}
	if count := handled.Load(); count != 18 {
		t.Errorf("expected 18 concurrently handled issues, received %d", count)
	}
	if count := sourceCalls.Load(); count != workers {
		t.Errorf("expected field source to run %d times, received %d", workers, count)
	}
}

func TestFieldSourceSupportsFutureDeclarationMechanisms(t *testing.T) {
	type payload struct {
		Name  string `future:"required"`
		Email string `future:"email"`
		Age   int
	}

	source := govalid.FieldSourceFunc(func(structType reflect.Type) ([]govalid.FieldSpec, error) {
		fields := make([]govalid.FieldSpec, 0)
		for index := 0; index < structType.NumField(); index++ {
			field := structType.Field(index)
			switch field.Tag.Get("future") {
			case "required":
				fields = append(fields, govalid.Field(field.Name, govalid.Required()))
			case "email":
				fields = append(fields, govalid.Field(field.Name, govalid.StringEmail()))
			}
		}
		return fields, nil
	})

	validator := govalid.New(govalid.WithFieldSources(source))
	err := validator.Validate(
		payload{Email: "invalid", Age: 10},
		govalid.Field("Age", govalid.IntMin(18)),
	)

	var validationErr *govalid.FieldIssueError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected FieldIssueError, received %T", err)
	}
	if len(validationErr.RulesIssues) != 3 {
		t.Fatalf("expected 2 discovered issues and 1 explicit issue, received %d", len(validationErr.RulesIssues))
	}
	if validationErr.RulesIssues[0].Path != "Name" ||
		validationErr.RulesIssues[1].Path != "Email" ||
		validationErr.RulesIssues[2].Path != "Age" {
		t.Fatalf("unexpected source and explicit field order: %#v", validationErr.RulesIssues)
	}
}

func TestFieldSourceErrorsAreReturned(t *testing.T) {
	type payload struct{ Name string }
	expected := errors.New("cannot discover fields")

	validator := govalid.New(govalid.WithFieldSources(
		govalid.FieldSourceFunc(func(reflect.Type) ([]govalid.FieldSpec, error) {
			return nil, expected
		}),
	))

	err := validator.Validate(payload{})
	if !errors.Is(err, expected) {
		t.Fatalf("expected source error to remain discoverable, received %v", err)
	}
}

func TestFieldIfPresentSupportsNilNestedParents(t *testing.T) {
	type contact struct {
		Email string
	}
	type payload struct {
		Contact *contact
	}

	validator := govalid.New()
	field := govalid.FieldIfPresent("Contact.Email", govalid.StringEmail())

	if err := validator.Validate(payload{}, field); err != nil {
		t.Fatalf("nil parent should skip optional nested field: %v", err)
	}
	if err := validator.Validate(payload{Contact: &contact{Email: "invalid"}}, field); err == nil {
		t.Fatal("present nested field should still be validated")
	}
	if err := validator.Validate(payload{}, govalid.Field("Contact.Email", govalid.StringEmail())); err == nil {
		t.Fatal("strict Field should report a nil intermediate path")
	}
}
