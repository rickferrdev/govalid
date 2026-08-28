package govalid

type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// IntPositive returns an integer validation rule for positive.
func IntPositive() Rule { return IntGreaterThan(0) }

// IntNegative returns an integer validation rule for negative.
func IntNegative() Rule { return IntLessThan(0) }

// IntNonPositive returns an integer validation rule for non positive.
func IntNonPositive() Rule { return IntMax(0) }

// IntNonNegative returns an integer validation rule for non negative.
func IntNonNegative() Rule { return IntMin(0) }

// IntZero returns an integer validation rule for zero.
func IntZero() Rule { return IntEqual(0) }

// IntNonZero returns an integer validation rule for non zero.
func IntNonZero() Rule { return IntNotEqual(0) }
