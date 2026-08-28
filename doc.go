// Package govalid provides composable, type-aware validation rules for
// selected fields of Go structs.
//
// Validation is explicit: callers associate field paths with rules through
// Field and execute them with a Validator. Rules cover strings, booleans,
// integers, floating-point numbers, bytes, maps, collections, universal value
// states, and conditional composition. Consumers can also implement custom
// rules with RuleContext and plug alternate field-declaration mechanisms into
// Validator with FieldSource.
package govalid
