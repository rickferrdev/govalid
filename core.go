package govalid

import (
	"errors"
	"reflect"
)

type Structurer struct {
	options options
}

type ruleContext struct {
	Root  reflect.Value
	Path  string
	Value reflect.Value
}

type Rule func(context ruleContext) error

type FieldSpec struct {
	Path  string
	Rules []Rule
}

type Issue struct {
	Path    string
	Message string
}

type FieldIssueError struct {
	RulesIssues []Issue
}

func New(options ...Option) *Structurer {
	config := defaultOptionsMap()

	for _, option := range options {
		if option != nil {
			option(&config)
		}
	}

	return &Structurer{
		options: config,
	}
}

func Field(path string, rules ...Rule) FieldSpec {
	return FieldSpec{
		Path:  path,
		Rules: rules,
	}
}

func (sttr *Structurer) Validate(data any, fields ...FieldSpec) error {
	root := reflect.ValueOf(data)

	if !root.IsValid() {
		return errors.New("a non-null structure was expected")
	}

	root = indirect(root)

	if root.Kind() != reflect.Struct {
		return errors.New("a struct value was expected")
	}

	var issues = make([]Issue, 0)
	for _, field := range fields {
		value, err := lookup(root, field.Path)
		if err != nil {
			return err
		}

		context := ruleContext{
			Root:  root,
			Path:  field.Path,
			Value: value,
		}

		for _, rule := range field.Rules {
			if err := rule(context); err != nil {
				issues = append(issues, Issue{
					Path:    field.Path,
					Message: err.Error(),
				})

				if sttr.options.StopOnFirstError {
					return &FieldIssueError{
						RulesIssues: issues,
					}
				}
			}
		}
	}

	if len(issues) == 0 {
		return nil
	}

	return &FieldIssueError{RulesIssues: issues}
}

func (issue *FieldIssueError) Error() string {
	return "validation failed"
}
