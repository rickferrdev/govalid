package govalid

import (
	"errors"
	"fmt"
)

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

// CollectionAll is an alias for CollectionEach.
func CollectionAll(rules ...Rule) Rule { return CollectionEach(rules...) }

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

func CollectionItemAt(index int, rules ...Rule) Rule {
	return collectionItemAt(index, false, rules)
}

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
