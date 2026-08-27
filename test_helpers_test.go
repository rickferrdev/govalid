package govalid

import "reflect"

func runRule(value any, rule Rule) error {
	return rule(ruleContext{Path: "Value", Value: reflect.ValueOf(value)})
}
