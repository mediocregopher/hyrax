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
	MsgType int
	Payload interface{}
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
func AddNode(addr *string) error {
	if err := ghost.AddConn(addr); err != nil {
		return err
	}
	ghost.SendAll(msg{MESH_ADD, *addr})
	return nil
}

// RemNode manually removes a node from this node's view of the mesh, and tells
// other nodes to do the same
func RemNode(addr *string) {
	ghost.SendAll(msg{MESH_REM, *addr})
	ghost.RemConn(addr)
}

// SendDirect sends a message directly to a node in the mesh
func SendDirect(addr string, m interface{}) error {
	return ghost.Send(&addr, msg{MESH_OTHER, m})
}

// SendAll sends a message to all the nodes currently connected to in the mesh
func SendAll(m interface{}) {
	ghost.SendAll(msg{MESH_OTHER, m})
}

func (ml *meshListen) msgSpin() {
	for m := range ml.inCh {
		switch m.(msg).MsgType {
		case MESH_ADD:
			addr := m.(msg).Payload.(string)
			ghost.AddConn(&addr)
		case MESH_REM:
			// TODO if we remove ourselves then we should disconnect from
			// everyone
			addr := m.(msg).Payload.(string)
			ghost.RemConn(&addr)
		case MESH_OTHER:
			ml.outCh <- m.(msg).Payload
		}
	}
}
