package govalid

type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

func IntPositive() Rule    { return IntGreaterThan(0) }
func IntNegative() Rule    { return IntLessThan(0) }
func IntNonPositive() Rule { return IntMax(0) }
func IntNonNegative() Rule { return IntMin(0) }
func IntZero() Rule        { return IntEqual(0) }
func IntNonZero() Rule     { return IntNotEqual(0) }
