package translate

import (
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

	// ToClientReturn takes in some bytes and decodes them into a ClientCommand,
	// or returns an error
	ToClientReturn([]byte) (*types.ClientReturn, error)

	// FromClientReturn takes in a client return and encodes it into a byte
	// slice, or returns an error if it can't
	FromClientReturn(*types.ClientReturn) ([]byte, error)

}
