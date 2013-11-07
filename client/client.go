package client

import (
	"errors"
	"github.com/mediocregopher/hyrax/client/net"
	"github.com/mediocregopher/hyrax/types"
	"github.com/mediocregopher/hyrax/translate"
)

// Client is an interface to interact with a connection to hyrax
type Client interface {
	
	// Cmd sends a command to hyrax and retrieves the result of the command,
	// either the return value or an error
	Cmd(*types.ClientCommand) (interface{}, error)

	// Close closes any connection the client may have with hyrax
	Close()

}

// NewClient takes in a format (ex. json), a connection type (ex. tcp) and an
// address where a hyrax server can be found (including the port). It also takes
// in a push channel, which can be nil if you want to ignore push messages. It
// returns a Client created from your specifications, or an error.
func NewClient(
	format, conntype, addr string,
	pushCh chan *types.ClientCommand) (Client, error) {

	trans, err := translate.StringToTranslator(format)
	if err != nil {
		return nil, err
	}

	switch conntype {
	case "tcp":
		return net.NewTcpClient(trans, addr, pushCh)
	default:
		return nil, errors.New("unknown connection type")
	}
}
