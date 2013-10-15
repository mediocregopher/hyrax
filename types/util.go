package types

// Byter is implemented by any value that has a Bytes method, which should
// return a byte-slice representation of the value. This is similar to
// fmt.Stringer.
type Byter interface {
	Bytes() []byte
}

// A simple wrapper for a byte slice to make it easy to work with other Byters
type SimpleByter []byte

func (b SimpleByter) Bytes() []byte {
	return []byte(b)
}

// Returns a SimpleByter. Doesn't actually do a copy.
func NewByter(b []byte) Byter {
	return SimpleByter(b)
}

// Given a separator and any number of byters, joins them into one string
func ByterJoin(sep Byter, b ...Byter) Byter {
	bb := make([][]byte, len(b))
	for i := range b {
		bb[i] = b[i].Bytes()
	}
	sepb := sep.Bytes()
	l := len(bb[0])
	for i := range bb[1:] {
		l += len(bb[i])
	}
	r := make([]byte, 0, l)
	r = append(r, bb[0]...)
	for i := range bb[1:] {
		r = append(r, sepb...)
		r = append(r, bb[i]...)
	}

	return NewByter(r)
}

// Uint64er is implemented by any value that has a Uint64 method, which should
// return a uint64 representation of the value.
type Uint64er interface {
	Uint64() uint64
}
