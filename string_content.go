package govalid

import (
	"fmt"
	"slices"
	"strings"
	"unicode"
)

func StringAlpha() Rule {
	return stringRuneRule("should contain only letters", unicode.IsLetter)
}

func StringAlphaNumeric() Rule {
	return stringRuneRule("should contain only letters and numbers", func(value rune) bool {
		return unicode.IsLetter(value) || unicode.IsNumber(value)
	})
}

func StringNumeric() Rule {
	return stringRuneRule("should contain only numbers", unicode.IsNumber)
}

func StringASCII() Rule {
	return stringRuneRule("should contain only ASCII characters", func(value rune) bool {
		return value <= unicode.MaxASCII
	})
}

func StringLowercase() Rule {
	return stringPredicateRule(
		func(value string) bool { return value == strings.ToLower(value) },
		"should be lowercase",
	)
}

func StringUppercase() Rule {
	return stringPredicateRule(
		func(value string) bool { return value == strings.ToUpper(value) },
		"should be uppercase",
	)
}

func StringTrimmed() Rule {
	return stringPredicateRule(
		func(value string) bool { return value == strings.TrimSpace(value) },
		"should not have leading or trailing spaces",
	)
}

func StringStartsWith(expect string) Rule {
	return stringPredicateRule(
		func(value string) bool { return strings.HasPrefix(value, expect) },
		fmt.Sprintf("should start with %s", expect),
	)
}

func StringStartsWithFold(expect string) Rule {
	return stringPredicateRule(
		func(value string) bool {
			return strings.HasPrefix(strings.ToLower(value), strings.ToLower(expect))
		},
		fmt.Sprintf("should start with %s, ignoring case", expect),
	)
}

func StringEndsWith(expect string) Rule {
	return stringPredicateRule(
		func(value string) bool { return strings.HasSuffix(value, expect) },
		fmt.Sprintf("should end with %s", expect),
	)
}

func StringEndsWithFold(expect string) Rule {
	return stringPredicateRule(
		func(value string) bool {
			return strings.HasSuffix(strings.ToLower(value), strings.ToLower(expect))
		},
		fmt.Sprintf("should end with %s, ignoring case", expect),
	)
}

func StringContains(expect string) Rule {
	return stringPredicateRule(
		func(value string) bool { return strings.Contains(value, expect) },
		fmt.Sprintf("should contain %s", expect),
	)
}

func StringContainsFold(expect string) Rule {
	return stringPredicateRule(
		func(value string) bool {
			return strings.Contains(strings.ToLower(value), strings.ToLower(expect))
		},
		fmt.Sprintf("should contain %s, ignoring case", expect),
	)
}

func StringNotContains(expect string) Rule {
	return stringPredicateRule(
		func(value string) bool { return !strings.Contains(value, expect) },
		fmt.Sprintf("should not contain %s", expect),
	)
}

func StringNotContainsFold(expect string) Rule {
	return stringPredicateRule(
		func(value string) bool {
			return !strings.Contains(strings.ToLower(value), strings.ToLower(expect))
		},
		fmt.Sprintf("should not contain %s, ignoring case", expect),
	)
}

func StringOneOf(expected ...string) Rule {
	return stringPredicateRule(
		func(value string) bool { return slices.Contains(expected, value) },
		fmt.Sprintf("should be one of %q", expected),
	)
}

func StringNotOneOf(unexpected ...string) Rule {
	return stringPredicateRule(
		func(value string) bool { return !slices.Contains(unexpected, value) },
		fmt.Sprintf("should not be one of %q", unexpected),
	)
}
