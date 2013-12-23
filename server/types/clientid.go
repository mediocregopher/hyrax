package types

import (
	"strconv"
)

// ClientId is a unique value that's given to every client of this hyrax node
type ClientId uint64

// Implements Bytes for the Byter interface
func (cid ClientId) Bytes() []byte {
	return []byte(strconv.FormatUint(uint64(cid), 16))
}

// Implements Uint64 for the Uint64er interface
func (cid ClientId) Uint64() uint64 {
	return uint64(cid)
}

// Given a byte slice (presumably returned from calling Bytes on a ClientId),
// returns the corresponding ClientId, and possibly an error if shit is fucked
func ClientIdFromBytes(b []byte) (ClientId, error) {
	n, err := strconv.ParseUint(string(b), 16, 64)
	if err != nil {
		return 0, err
	}

	return ClientId(n), nil
}

// Given a uint64 (presumably returned from calling Uint64() on a ClientId),
// returns the corresponding ClientId, and possibly an error if shit is fucked
func ClientIdFromUint64(i uint64) (ClientId, error) {
	return ClientId(i), nil
}
