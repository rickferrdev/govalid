package govalid

import (
	"errors"
	"fmt"
)

// CollectionEach returns a collection validation rule for each.
func CollectionEach(rules ...Rule) Rule {
	return func(context ruleContext) error {
		value, err := collectionValue(context)
		if err != nil {
			return err
		}
		for index := 0; index < value.Len(); index++ {
			itemContext := collectionItemContext(context, value.Index(index), index)
			if err := applyRules(itemContext, rules); err != nil {
				return fmt.Errorf("invalid collection item at %s: %w", itemContext.Path, err)
			}
		}
		return nil
	}
}

// CollectionAny returns a collection validation rule for any.
func CollectionAny(rules ...Rule) Rule {
	return func(context ruleContext) error {
		value, err := collectionValue(context)
		if err != nil {
			return err
		}
		for index := 0; index < value.Len(); index++ {
			if applyRules(collectionItemContext(context, value.Index(index), index), rules) == nil {
				return nil
			}
		}
		return errors.New("at least one collection item should satisfy all rules")
	}
}

// CollectionNone returns a collection validation rule for none.
func CollectionNone(rules ...Rule) Rule {
	return func(context ruleContext) error {
		value, err := collectionValue(context)
		if err != nil {
			return err
		}
		for index := 0; index < value.Len(); index++ {
			if applyRules(collectionItemContext(context, value.Index(index), index), rules) == nil {
				return fmt.Errorf("collection item at index %d should not satisfy all rules", index)
			}
		}
		return nil
	}
}

// CollectionItemAt returns a collection validation rule for item at.
func CollectionItemAt(index int, rules ...Rule) Rule {
	return collectionItemAt(index, false, rules)
}

// CollectionItemAtIfPresent returns a collection validation rule for item at if present.
func CollectionItemAtIfPresent(index int, rules ...Rule) Rule {
	return collectionItemAt(index, true, rules)
}

func collectionItemAt(index int, optional bool, rules []Rule) Rule {
	return func(context ruleContext) error {
		value, err := collectionValue(context)
		if err != nil {
			return err
		}
		if index < 0 {
			return errors.New("collection index cannot be negative")
		}
		if index >= value.Len() {
			if optional {
				return nil
			}
			return fmt.Errorf("collection should contain index %d", index)
		}
		itemContext := collectionItemContext(context, value.Index(index), index)
		if err := applyRules(itemContext, rules); err != nil {
			return fmt.Errorf("invalid collection item at %s: %w", itemContext.Path, err)
		}
		return nil
	}
}
