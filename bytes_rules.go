package govalid

import (
	"errors"
	"fmt"
)

func BytesEach(rules ...Rule) Rule {
	return func(context ruleContext) error {
		value, err := bytesValue(context)
		if err != nil {
			return err
		}
		for index := 0; index < value.Len(); index++ {
			itemContext := ruleContext{
				Root: context.Root, Path: fmt.Sprintf("%s[%d]", context.Path, index), Value: value.Index(index),
			}
			if err := applyRules(itemContext, rules); err != nil {
				return fmt.Errorf("invalid byte at %s: %w", itemContext.Path, err)
			}
		}
		return nil
	}
}

func BytesAt(index int, rules ...Rule) Rule {
	return bytesAt(index, false, rules)
}

func BytesAtIfPresent(index int, rules ...Rule) Rule {
	return bytesAt(index, true, rules)
}

func bytesAt(index int, optional bool, rules []Rule) Rule {
	return func(context ruleContext) error {
		value, err := bytesValue(context)
		if err != nil {
			return err
		}
		if index < 0 {
			return errors.New("byte index cannot be negative")
		}
		if index >= value.Len() {
			if optional {
				return nil
			}
			return fmt.Errorf("bytes should contain index %d", index)
		}
		itemContext := ruleContext{
			Root: context.Root, Path: fmt.Sprintf("%s[%d]", context.Path, index), Value: value.Index(index),
		}
		if err := applyRules(itemContext, rules); err != nil {
			return fmt.Errorf("invalid byte at %s: %w", itemContext.Path, err)
		}
		return nil
	}
}
