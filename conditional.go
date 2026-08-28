package govalid

import "errors"

// Condition decides whether contextual rules should run.
type Condition func(context RuleContext) bool

// When applies rules only when condition is true.
func When(condition bool, rules ...Rule) Rule {
	return func(context RuleContext) error {
		if !condition {
			return nil
		}
		return applyRules(context, rules)
	}
}

// Unless applies rules only when condition is false.
func Unless(condition bool, rules ...Rule) Rule {
	return When(!condition, rules...)
}

// WhenContext applies rules when condition accepts the current RuleContext.
func WhenContext(condition Condition, rules ...Rule) Rule {
	return func(context RuleContext) error {
		if condition == nil {
			return errors.New("condition should not be nil")
		}
		if !condition(context) {
			return nil
		}
		return applyRules(context, rules)
	}
}

// UnlessContext applies rules when condition rejects the current RuleContext.
func UnlessContext(condition Condition, rules ...Rule) Rule {
	return func(context RuleContext) error {
		if condition == nil {
			return errors.New("condition should not be nil")
		}
		if condition(context) {
			return nil
		}
		return applyRules(context, rules)
	}
}

// Optional skips rules for nil and empty values. Scalar values, including 0
// and false, are not considered absent and are validated normally.
func Optional(rules ...Rule) Rule {
	return func(context RuleContext) error {
		if isAbsentValue(context.Value) {
			return nil
		}
		return applyRules(context, rules)
	}
}
