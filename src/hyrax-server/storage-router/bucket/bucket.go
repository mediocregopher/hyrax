package bucket

import (
	"sync"
	"fmt"
	"github.com/mediocregopher/hyrax/src/hyrax/types"
)

var bList []*string
var bLock sync.RWMutex

// Set sets the bucket at bIndex to the given string. The index must be inside
// the existing bucket list, or be exactly one position outside of it. In
// effect, this method is SetOrAppend
func Set(s *string, bIndex int) error {
	bLock.Lock()
	defer bLock.Unlock()

	bLen := len(bList)
	if bIndex > bLen {
		return fmt.Errorf("The bucket list length is %d, its size can only be increased by one, so it's maximum index right now is %d", bLen, bLen)
	} else if bIndex == bLen {
		bList = append(bList,s)
	} else {
		bList[bIndex] = s
	}
	return nil
}

// Buckets returns a copy (not the actual slice, so you don't have to worry
// about it being changed under your nose) of the bucket list in its current
// state
func Buckets() []*string {
	bLock.RLock()
	defer bLock.RUnlock()

	r := make([]*string, len(bList))
	for i := range bList {
		r[i] = bList[i]
	}
	return r
}

// KeyBucket hashes the given key with a Shift-Add-XOR hash, identifies which
// bucket it belongs in based on this hash, and returns the string at that
// bucket location.
func KeyBucket(b types.Byter) (*string, error) {
	bb := b.Bytes()
	h := uint(0)
	for _, l := range bb {
		h = h ^ ((h << 5) + (h >> 2) + uint(l))
	}
	bLock.RLock()
	defer bLock.RUnlock()

	if bList == nil {
		return nil, fmt.Errorf("No strings set in the bucket")
	}

	pos := h % uint(len(bList))
	if s := bList[pos]; s != nil {
		return s, nil
	} else {
		return nil, fmt.Errorf("There is no string set at bucket pos %d", pos)
	}
}
