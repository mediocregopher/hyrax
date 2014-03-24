package types

import (
	"github.com/mediocregopher/hyrax/types"
)

// Client is an interface which must be implemented by clients to hyrax (go
// figure)
type Client interface {

	// ClientId returns the ClientId of a given client (again, go figure)
	ClientId() ClientId

	// CommandPushCh returns a channel where commands that are to be pushed
	// to the client should be pushed on to
	CommandPushCh() chan<- *types.ClientCommand

	// ClosingCh returns a channel which will have close() called on it when the
	// connection is closed
	ClosingCh() <-chan struct{}
}
