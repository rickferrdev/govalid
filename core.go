package govalid

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// Validator validates selected fields of Go structs.
// A Validator is safe for concurrent use after construction when its issue
// handler, if any, is also concurrency-safe.
type Validator struct {
	options options
}

// RuleContext contains the values available to a validation rule.
//
// Root is the validated struct, Value is the selected field, and Path is the
// field path passed to Field.
type RuleContext struct {
	Root  reflect.Value
	Path  string
	Value reflect.Value
}

// ValueAny returns the selected field as an interface value. It returns nil
// when the value is invalid or cannot be interfaced.
func (context RuleContext) ValueAny() any {
	if !context.Value.IsValid() || !context.Value.CanInterface() {
		return nil
	}
	return context.Value.Interface()
}

// RootAny returns the validated struct as an interface value. It returns nil
// when the root is invalid or cannot be interfaced.
func (context RuleContext) RootAny() any {
	if !context.Root.IsValid() || !context.Root.CanInterface() {
		return nil
	}
	return context.Root.Interface()
}

// ruleContext is kept internally so existing rule implementations remain
// concise while RuleContext is available to consumers.
type ruleContext = RuleContext

// Rule validates the field described by a RuleContext.
type Rule func(context RuleContext) error

// FieldSpec associates a struct field path with validation rules.
type FieldSpec struct {
	Path          string
	Rules         []Rule
	skipNilParent bool
}

// FieldSource discovers field specifications for a struct type. It allows
// alternate declaration mechanisms, such as a future struct-tag parser, to
// reuse Validator without changing its execution API.
//
// Implementations used concurrently must be concurrency-safe.
type FieldSource interface {
	Fields(structType reflect.Type) ([]FieldSpec, error)
}

// FieldSourceFunc adapts a function to FieldSource.
type FieldSourceFunc func(structType reflect.Type) ([]FieldSpec, error)

// Fields discovers field specifications for structType.
func (source FieldSourceFunc) Fields(structType reflect.Type) ([]FieldSpec, error) {
	if source == nil {
		return nil, errors.New("field source function should not be nil")
	}
	return source(structType)
}

// Issue describes one failed validation rule.
type Issue struct {
	Path    string
	Value   any
	Rule    string
	Message string
}

// FieldIssueError contains the issues collected during one validation run.
type FieldIssueError struct {
	RulesIssues []Issue
}

// As lets errors.As extract the first individual Issue from a validation
// error. Extract FieldIssueError itself to inspect every collected issue.
func (field *FieldIssueError) As(target any) bool {
	issueTarget, ok := target.(**Issue)
	if !ok || len(field.RulesIssues) == 0 {
		return false
	}

	issue := field.RulesIssues[0]
	*issueTarget = &issue
	return true
}

// Error formats every collected validation issue.
func (field *FieldIssueError) Error() string {
	var builder strings.Builder
	builder.WriteString("validation error:\n")

	for index, issue := range field.RulesIssues {
		fmt.Fprintf(&builder, "(%d)\trule: %s\n\tpath: %s\n\tvalue: %s\n\tmessage: %s\n",
			index,
			issue.Rule,
			issue.Path,
			formatIssueValue(issue.Value),
			issue.Message,
		)
	}

	return builder.String()
}

// Error formats the failed rule, path, rejected value, and message.
func (issue *Issue) Error() string {
	var builder strings.Builder
	builder.WriteString("issue error: ")

	fmt.Fprintf(&builder, "rule: %s path: %s value: %s message: %s\n",
		issue.Rule,
		issue.Path,
		formatIssueValue(issue.Value),
		issue.Message,
	)

	return builder.String()
}

func formatIssueValue(value any) string {
	if value == nil {
		return "nil (<nil>)"
	}

	reflected := reflect.ValueOf(value)
	switch reflected.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice:
		if reflected.Len() > 32 {
			return fmt.Sprintf("<%s len=%d> (%T)", reflected.Type(), reflected.Len(), value)
		}
	case reflect.String:
		const maximumRunes = 256
		runes := []rune(reflected.String())
		if len(runes) > maximumRunes {
			return fmt.Sprintf("%q... (%T, %d runes)", string(runes[:maximumRunes]), value, len(runes))
		}
	}

	return fmt.Sprintf("%#v (%T)", value, value)
}

// New creates a Validator configured with options.
func New(options ...Option) *Validator {
	config := defaultOptionsMap()

	for _, option := range options {
		if option != nil {
			option(&config)
		}
	}

	return &Validator{
		options: config,
	}
}

// Field associates a dot-separated struct field path with validation rules.
func Field(path string, rules ...Rule) FieldSpec {
	return FieldSpec{
		Path:  path,
		Rules: rules,
	}
}

// FieldIfPresent associates rules with a field path but skips them when an
// intermediate pointer or interface in that path is nil. The selected field
// itself is still validated when its parent path exists.
func FieldIfPresent(path string, rules ...Rule) FieldSpec {
	return FieldSpec{
		Path:          path,
		Rules:         rules,
		skipNilParent: true,
	}
}

// Validate applies field specifications to a struct or pointer to a struct.
// It returns nil when every rule passes or when errors are explicitly silenced.
func (sttr *Validator) Validate(data any, fields ...FieldSpec) error {
	root := reflect.ValueOf(data)

	if !root.IsValid() {
		return errors.New("a non-null structure was expected")
	}

	root = indirect(root)

	if root.Kind() != reflect.Struct {
		return errors.New("a struct value was expected")
	}

	resolvedFields := make([]FieldSpec, 0, len(fields))
	for index, source := range sttr.options.FieldSources {
		if source == nil {
			continue
		}
		discovered, err := source.Fields(root.Type())
		if err != nil {
			return fmt.Errorf("field source %d: %w", index, err)
		}
		resolvedFields = append(resolvedFields, discovered...)
	}
	resolvedFields = append(resolvedFields, fields...)

	var issues = make([]Issue, 0)
	for _, field := range resolvedFields {
		value, err := lookup(root, field.Path)
		if err != nil {
			var nilParent *nilPathError
			if field.skipNilParent && errors.As(err, &nilParent) {
				continue
			}
			return err
		}

		context := ruleContext{
			Root:  root,
			Path:  field.Path,
			Value: value,
		}

		for _, rule := range field.Rules {
			if rule == nil {
				continue
			}
			if err := rule(context); err != nil {
				issue := Issue{
					Path:    field.Path,
					Rule:    ruleName(rule),
					Value:   ruleValue(value),
					Message: err.Error(),
				}
				issues = append(issues, issue)

				if sttr.options.IssueHandler != nil {
					sttr.options.IssueHandler(issue)
				}

				if sttr.options.WithPanicOnFirstError {
					panic(&issue)
				}

				if sttr.options.StopOnFirstError {
					if sttr.options.WithSilenceErrors {
						return nil
					}
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

	if sttr.options.WithSilenceErrors {
		return nil
	}

	return &FieldIssueError{RulesIssues: issues}
}

func ruleName(rule Rule) string {
	function := runtime.FuncForPC(reflect.ValueOf(rule).Pointer())
	if function == nil {
		return "unknown"
	}

	name := function.Name()
	if slash := strings.LastIndex(name, "/"); slash >= 0 {
		name = name[slash+1:]
	}
	return strings.TrimSuffix(name, "-fm")
}

func ruleValue(value reflect.Value) any {
	if !value.IsValid() {
		return nil
	}
	if value.CanInterface() {
		return value.Interface()
	}
	return fmt.Sprintf("<%s value>", value.Type())
}
