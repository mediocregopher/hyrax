package client

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"errors"

	"github.com/mediocregopher/hyrax/client/net"
	"github.com/mediocregopher/hyrax/translate"
	"github.com/mediocregopher/hyrax/types"
)

// Client is an interface to interact with a connection to hyrax
type Client interface {

	// Cmd sends a command to hyrax and retrieves the result of the command,
	// either the return value or an error. The error will be io.EOF if and only
	// if the connection has been closed
	Cmd(*types.ClientCommand) (interface{}, error)

	// Close closes any connection the client may have with hyrax
	Close()
}

// NewClient takes in a format (ex. json), a connection type (ex. tcp) and an
// address where a hyrax server can be found (including the port). It also takes
// in a push channel, which can be nil if you want to ignore push messages. It
// returns a Client created from your specifications, or an error.
func NewClient(
	le *types.ListenEndpoint, pushCh chan *types.ClientCommand) (Client, error) {

	trans, err := translate.StringToTranslator(le.Format)
	if err != nil {
		return nil, err
	}

	switch le.Type {
	case "tcp":
		return net.NewTcpClient(trans, le.Addr, pushCh)
	default:
		return nil, errors.New("unknown connection type")
	}
}

// Given a command and a secret used to generate the hash for a command, does
// all the work of actually creating a ClientCommand
func CreateClientCommand(
	cmd, keyB, id, secretKey []byte,
	args ...interface{}) *types.ClientCommand {

	mac := hmac.New(sha1.New, secretKey)
	mac.Write(cmd)
	mac.Write(keyB)
	mac.Write(id)
	sum := mac.Sum(nil)
	sumhex := make([]byte, hex.EncodedLen(len(sum)))
	hex.Encode(sumhex, sum)

	return &types.ClientCommand{
		Command:    cmd,
		StorageKey: keyB,
		Args:       args,
		Id:         id,
		Secret:     sumhex,
	}
}
