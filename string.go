package govalid

import "errors"

func StringRequired() Rule {
	return func(context ruleContext) error {
		value, err := stringValue(context)
		if err != nil {
			return err
		}
		if value == "" {
			return errors.New("should not be empty")
		}
		return nil
	}
}
