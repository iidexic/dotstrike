package meta

type State interface {
	Trigger() error
	Next() State
}

type Condition interface {
	Met() bool
}
type Settable[A any] interface {
	Set(A ...A) error
}

type SelfComparable[C comparable] interface {
	Equal(C) bool
}

type Equals[C comparable] struct {
	basis C
}

type cmpint int

func (c cmpint) Equal(c2 cmpint) bool { return c == c2 }

// booleans, numbers, strings, pointers, channels, arrays of comparable types,
// structs whose fields are all comparable types
func (E *Equals[C]) Set(c1 C) bool {
	return E.basis == c1
}
