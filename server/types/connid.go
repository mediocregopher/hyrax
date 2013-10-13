package types

import (
	"encoding/binary"
)

// ConnId is a unique value that's given to every connection on this hyrax node
type ConnId uint64

// Implements Bytes for the Byter interface
func (cid *ConnId) Bytes() []byte {
	size := binary.Size(*cid)
	buf := make([]byte, size)
	binary.PutUvarint(buf, uint64(*cid))
	return buf
}

// Implements Uint64 for the Uint64er interface
func (cid *ConnId) Uint64() uint64 {
	return uint64(*cid)
}
