package types

import (
	"encoding/binary"
	"fmt"
)

// ClientId is a unique value that's given to every client of this hyrax node
type ClientId uint64

// Implements Bytes for the Byter interface
func (cid ClientId) Bytes() []byte {
	size := binary.Size(cid)
	buf := make([]byte, size)
	binary.PutUvarint(buf, uint64(cid))
	return buf
}

// Implements Uint64 for the Uint64er interface
func (cid ClientId) Uint64() uint64 {
	return uint64(cid)
}

// Given a byte slice (presumably returned from calling Bytes on a ClientId),
// returns the corresponding ClientId, and possibly an error if shit is fucked
func ClientIdFromBytes(b []byte) (ClientId, error) {
	cid, i := binary.Uvarint(b)
	if i <= 0 {
		return 0, fmt.Errorf("Error reading ClientId from bytes: %v", b)
	}
	return ClientId(cid), nil
}

// Given a uint64 (presumably returned from calling Uint64() on a ClientId),
// returns the corresponding ClientId, and possibly an error if shit is fucked
func ClientIdFromUint64(i uint64) (ClientId, error) {
	return ClientId(i), nil
}
