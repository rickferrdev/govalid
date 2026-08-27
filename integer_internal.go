package govalid

import (
	"cmp"
	"errors"
	"reflect"
	"strconv"
)

type number struct {
	signed bool
	i64    int64
	u64    uint64
}

func numberFrom[T integer](value T) number {
	rv := reflect.ValueOf(value)
	if rv.Kind() >= reflect.Int && rv.Kind() <= reflect.Int64 {
		return number{signed: true, i64: rv.Int()}
	}
	return number{u64: rv.Uint()}
}

func integerValue(context ruleContext) (number, error) {
	switch context.Value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return number{signed: true, i64: context.Value.Int()}, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return number{u64: context.Value.Uint()}, nil
	default:
		return number{}, errors.New("should be of integer type")
	}
}

func integerRule(validate func(number) error) Rule {
	return func(context ruleContext) error {
		value, err := integerValue(context)
		if err != nil {
			return err
		}
		return validate(value)
	}
}

func (n number) compare(other number) int {
	if n.signed && other.signed {
		return cmp.Compare(n.i64, other.i64)
	}
	if !n.signed && !other.signed {
		return cmp.Compare(n.u64, other.u64)
	}
	if n.signed {
		if n.i64 < 0 {
			return -1
		}
		return cmp.Compare(uint64(n.i64), other.u64)
	}
	if other.i64 < 0 {
		return 1
	}
	return cmp.Compare(n.u64, uint64(other.i64))
}

func (n number) isNegative() bool { return n.signed && n.i64 < 0 }

func (n number) magnitude() uint64 {
	if !n.signed {
		return n.u64
	}
	if n.i64 >= 0 {
		return uint64(n.i64)
	}
	return uint64(-(n.i64 + 1)) + 1
}

func (n number) String() string {
	if n.signed {
		return strconv.FormatInt(n.i64, 10)
	}
	return strconv.FormatUint(n.u64, 10)
}
