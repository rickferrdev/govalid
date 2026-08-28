package govalid

import (
	"fmt"
	"slices"
	"strings"
	"unicode"
)

// StringAlpha returns a string validation rule for alpha.
func StringAlpha() Rule {
	return stringRuneRule("should contain only letters", unicode.IsLetter)
}

// StringAlphaNumeric returns a string validation rule for alpha numeric.
func StringAlphaNumeric() Rule {
	return stringRuneRule("should contain only letters and numbers", func(value rune) bool {
		return unicode.IsLetter(value) || unicode.IsNumber(value)
	})
}

// StringNumeric returns a string validation rule for numeric.
func StringNumeric() Rule {
	return stringRuneRule("should contain only numbers", unicode.IsNumber)
}

// StringASCII returns a string validation rule for ascii.
func StringASCII() Rule {
	return stringRuneRule("should contain only ASCII characters", func(value rune) bool {
		return value <= unicode.MaxASCII
	})
}

// StringLowercase returns a string validation rule for lowercase.
func StringLowercase() Rule {
	return stringPredicateRule(
		func(value string) bool { return value == strings.ToLower(value) },
		"should be lowercase",
	)
}

// StringUppercase returns a string validation rule for uppercase.
func StringUppercase() Rule {
	return stringPredicateRule(
		func(value string) bool { return value == strings.ToUpper(value) },
		"should be uppercase",
	)
}

// StringTrimmed returns a string validation rule for trimmed.
func StringTrimmed() Rule {
	return stringPredicateRule(
		func(value string) bool { return value == strings.TrimSpace(value) },
		"should not have leading or trailing spaces",
	)
}

// StringStartsWith returns a string validation rule for starts with.
func StringStartsWith(expect string) Rule {
	return stringPredicateRule(
		func(value string) bool { return strings.HasPrefix(value, expect) },
		fmt.Sprintf("should start with %s", expect),
	)
}

// StringStartsWithFold returns a string validation rule for starts with fold.
func StringStartsWithFold(expect string) Rule {
	return stringPredicateRule(
		func(value string) bool {
			return strings.HasPrefix(strings.ToLower(value), strings.ToLower(expect))
		},
		fmt.Sprintf("should start with %s, ignoring case", expect),
	)
}

// StringEndsWith returns a string validation rule for ends with.
func StringEndsWith(expect string) Rule {
	return stringPredicateRule(
		func(value string) bool { return strings.HasSuffix(value, expect) },
		fmt.Sprintf("should end with %s", expect),
	)
}

// StringEndsWithFold returns a string validation rule for ends with fold.
func StringEndsWithFold(expect string) Rule {
	return stringPredicateRule(
		func(value string) bool {
			return strings.HasSuffix(strings.ToLower(value), strings.ToLower(expect))
		},
		fmt.Sprintf("should end with %s, ignoring case", expect),
	)
}

// StringContains returns a string validation rule for contains.
func StringContains(expect string) Rule {
	return stringPredicateRule(
		func(value string) bool { return strings.Contains(value, expect) },
		fmt.Sprintf("should contain %s", expect),
	)
}

// StringContainsFold returns a string validation rule for contains fold.
func StringContainsFold(expect string) Rule {
	return stringPredicateRule(
		func(value string) bool {
			return strings.Contains(strings.ToLower(value), strings.ToLower(expect))
		},
		fmt.Sprintf("should contain %s, ignoring case", expect),
	)
}

// StringNotContains returns a string validation rule for not contains.
func StringNotContains(expect string) Rule {
	return stringPredicateRule(
		func(value string) bool { return !strings.Contains(value, expect) },
		fmt.Sprintf("should not contain %s", expect),
	)
}

// StringNotContainsFold returns a string validation rule for not contains fold.
func StringNotContainsFold(expect string) Rule {
	return stringPredicateRule(
		func(value string) bool {
			return !strings.Contains(strings.ToLower(value), strings.ToLower(expect))
		},
		fmt.Sprintf("should not contain %s, ignoring case", expect),
	)
}

// StringOneOf returns a string validation rule for one of.
func StringOneOf(expected ...string) Rule {
	return stringPredicateRule(
		func(value string) bool { return slices.Contains(expected, value) },
		fmt.Sprintf("should be one of %q", expected),
	)
}

// StringNotOneOf returns a string validation rule for not one of.
func StringNotOneOf(unexpected ...string) Rule {
	return stringPredicateRule(
		func(value string) bool { return !slices.Contains(unexpected, value) },
		fmt.Sprintf("should not be one of %q", unexpected),
	)
}
