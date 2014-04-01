package translate

import (
	"fmt"
	"strings"

	"github.com/mediocregopher/hyrax/types"
)

// A translator takes in raw data of some known format and translates it into
// types that hyrax will use for things
type Translator interface {

	// ToClientCommand takes in some bytes and decodes them into a
	// ClientCommand, or returns an error
	ToClientCommand([]byte) (*types.ClientCommand, error)

	// FromClientCommand takes in a client command and encodes it into a byte
	// slice, or returns an error if it can't
	FromClientCommand(*types.ClientCommand) ([]byte, error)

	// ToClientReturn takes in some bytes and decodes them into a ClientReturn,
	// or returns an error
	ToClientReturn([]byte) (*types.ClientReturn, error)

	// FromClientReturn takes in a client return and encodes it into a byte
	// slice, or returns an error if it can't
	FromClientReturn(*types.ClientReturn) ([]byte, error)
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
