package mesh

import (
	"github.com/mediocregopher/ghost"
)

const (
	MESH_ADD = iota
	MESH_REM
	MESH_OTHER
)

type msg struct {
	msgtype int
	payload interface{}
}

type meshListen struct {
	inCh chan interface{}
	outCh chan interface{}
	errCh chan error
}

func init() {
	ghost.Register(msg{})
}

// Listen begins listening for other mesh connections on the given address, and
// starts responding to add/remove node messages coming from the network. It
// returns a channel that other kinds of messages from the network will be
// pushed to, a channel that network errors will be pushed to, and an error if
// listening failed for some reason
func Listen(addr string) (<-chan interface{}, <-chan error, error) {
	msgCh, errCh, err := ghost.Listen(addr)
	if err != nil {
		return nil, nil, err
	}

	ml := meshListen{
		inCh: msgCh,
		outCh: make(chan interface{}),
		errCh: errCh,
	}
	go ml.msgSpin()

	return ml.outCh, ml.errCh, nil
}

// RegisterMsgType is used to register structs that will be being sent to other
// nodes, or received from them. If you don't do this then it won't be possible
// to decode your message
func RegisterMsgType(typ interface{}) {
	ghost.Register(typ)
}

// AddNode manually adds a node to this node's view of the mesh, and tells other
// nodes to do the same
func AddNode(addr string) {
	// TODO AddConn isn't synchronous. We want the new node to add itself to its
	// mesh view, but it's possible this won't happen
	ghost.AddConn(addr)
	ghost.SendAll(msg{MESH_ADD, addr})
}

// RemNode manually removes a node from this node's view of the mesh, and tells
// other nodes to do the same
func RemNode(addr string) {
	ghost.RemoveConn(addr)
	ghost.SendAll(msg{MESH_REM, addr})
}

// SendDirect sends a message directly to a node in the mesh
func SendDirect(addr string, m interface{}) {
	ghost.Send(addr, msg{MESH_OTHER, m})
}

// SendAll sends a message to all the nodes currently connected to in the mesh
func SendAll(m interface{}) {
	ghost.SendAll(msg{MESH_OTHER, m})
}

func (ml *meshListen) msgSpin() {
	for m := range ml.inCh {
		switch m.(msg).msgtype {
		case MESH_ADD:
			ghost.AddConn(m.(msg).payload.(string))
		case MESH_REM:
			// TODO if we remove ourselves then we should disconnect from
			// everyone
			ghost.RemoveConn(m.(msg).payload.(string))
		case MESH_OTHER:
			ml.outCh <- m.(msg).payload
		}
	}
}
