package types

import (
	"github.com/mediocregopher/hyrax/types"
)

// Client is an interface which must be implemented by clients to hyrax (go
// figure)
type Client interface {

	// ClientId returns the ClientId of a given client (again, go figure)
	ClientId() ClientId

	// PushCh returns a channel where actions that are to be pushed to the
	// client should be pushed on to
	PushCh() chan<- *types.Action

	// ClosingCh returns a channel which will have close() called on it when the
	// connection is closed
	ClosingCh() <-chan struct{}
}
