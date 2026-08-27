package govalid

import (
	"fmt"
	"unicode/utf8"
)

func StringLength(expect int) Rule {
	return stringLengthRule(
		func(length int) bool { return length == expect },
		fmt.Sprintf("should be %d in length", expect),
	)
}

func StringMinLength(expect int) Rule {
	return stringLengthRule(
		func(length int) bool { return length >= expect },
		fmt.Sprintf("should have a minimum size of %d", expect),
	)
}

func StringMaxLength(expect int) Rule {
	return stringLengthRule(
		func(length int) bool { return length <= expect },
		fmt.Sprintf("should have a maximum size of %d", expect),
	)
}

func stringLengthRule(valid func(int) bool, message string) Rule {
	return stringPredicateRule(
		func(value string) bool { return valid(utf8.RuneCountInString(value)) },
		message,
	)
}
