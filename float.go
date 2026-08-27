package govalid

type floating interface {
	~float32 | ~float64
}

type floatNumber struct {
	value float64
	bits  int
}
