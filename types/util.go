package types

// Byter is implemented by any value that has a Bytes method, which should
// return a byte-slice representation of the value. This is similar to
// fmt.Stringer.
type Byter interface {
	Bytes() []byte
}

// Uint64er is implemented by any value that has a Uint64 method, which should
// return a uint64 representation of the value.
type Uint64er interface {
	Uint64() uint64
}
