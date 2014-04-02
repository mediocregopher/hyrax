package translate

import (
	"fmt"
	"strings"

	"github.com/mediocregopher/hyrax/types"
)

// A translator takes in raw data of some known format and translates it into
// types that hyrax will use for things
type Translator interface {

	// ToAction takes in some bytes and decodes them into a
	// Action, or returns an error
	ToAction([]byte) (*types.Action, error)

	// FromAction takes in a client command and encodes it into a byte
	// slice, or returns an error if it can't
	FromAction(*types.Action) ([]byte, error)

	// ToActionReturn takes in some bytes and decodes them into a ActionReturn,
	// or returns an error
	ToActionReturn([]byte) (*types.ActionReturn, error)

	// FromActionReturn takes in a client return and encodes it into a byte
	// slice, or returns an error if it can't
	FromActionReturn(*types.ActionReturn) ([]byte, error)
}

// StringToTranslator takes in a string which is supposed to identify which
// translator is desired and returns an instance of a translator of that type
func StringToTranslator(ts string) (Translator, error) {
	switch strings.ToLower(ts) {
	case "json":
		return &JsonTranslator{}, nil
	default:
		return nil, fmt.Errorf("unknown format: %s", ts)
	}
}
