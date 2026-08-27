package govalid

import (
	"errors"
	"net/mail"
	"net/url"
	"regexp"
)

var uuidPattern = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func StringRegex(expression regexp.Regexp) Rule {
	return stringPredicateRule(
		expression.MatchString,
		"should match the regex expression",
	)
}

func StringEmail() Rule {
	return func(context ruleContext) error {
		value, err := stringValue(context)
		if err != nil {
			return err
		}
		if value == "" {
			return nil
		}

		address, parseErr := mail.ParseAddress(value)
		if parseErr != nil || address.Name != "" || address.Address != value {
			return errors.New("should be a valid email address")
		}
		return nil
	}
}

func StringURL() Rule {
	return func(context ruleContext) error {
		value, err := stringValue(context)
		if err != nil {
			return err
		}
		if value == "" {
			return nil
		}

		parsed, parseErr := url.ParseRequestURI(value)
		if parseErr != nil || parsed.Scheme == "" || parsed.Host == "" {
			return errors.New("should be a valid absolute URL")
		}
		return nil
	}
}

func StringUUID() Rule {
	return func(context ruleContext) error {
		value, err := stringValue(context)
		if err != nil {
			return err
		}
		if value == "" {
			return nil
		}
		if !uuidPattern.MatchString(value) {
			return errors.New("should be a valid UUID")
		}
		return nil
	}
}
