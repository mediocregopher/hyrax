package dist

import (
	"github.com/mediocregopher/hyrax/server/dist/mesh"
	"github.com/mediocregopher/hyrax/types"
)

// StorageBucketSet is used to tell a node to add a storage unit to a particular
// bucket
type StorageBucketSet struct {
	Bucket         int
	ConnType, Addr string
	Extra          []interface{}
}

func init() {
	mesh.RegisterMsgType(StorageBucketSet{})
	mesh.RegisterMsgType(types.ClientCommand{})
}

type distWorker struct {
	// These two are provided by mesh, and are for incoming messages only
	msgCh <-chan interface{}
	errCh <-chan error

	clientCmdCh chan *types.ClientCommand
	setBucketCh chan *StorageBucketSet
}

var distWorkerInst distWorker

// A channel where all StorageBucketSet commands being sent by other nodes get
// routed to
var StorageBucketSetCommands = make(chan *StorageBucketSet)

// A channel where all ClientCommands being sent by other nodes get routed to
var ClientCommands = make(chan *types.ClientCommand)

// Initializes the dist worker, which will set up a mesh listener on the given
// address and pull and handle data from that. It will also pull and handle
// events that need sending to other nodes from all over hyrax
func Init(addr string) error {
	msgCh, errCh, err := mesh.Listen(addr)
	if err != nil {
		return err
	}

	distWorkerInst = distWorker{
		msgCh:       msgCh,
		errCh:       errCh,
		clientCmdCh: make(chan *types.ClientCommand),
		setBucketCh: make(chan *StorageBucketSet),
	}
	go distWorkerInst.spin()

	return nil
}

// AddNode passes straight through to mesh.AddNode
func AddNode(addr *string) error {
	return mesh.AddNode(addr)
}

// RemNode passes straight through to mesh.RemNode
func RemNode(addr *string) {
	mesh.RemNode(addr)
}

// Sends a client command to all nodes to be routed to clients monitoring a key
func SendClientCommand(cmd *types.ClientCommand) {
	distWorkerInst.clientCmdCh <- cmd
}

// Tells all nodes to add a storage bucket
func SendSetBucket(bIndex int, conntype, addr string, extra ...interface{}) {
	distWorkerInst.setBucketCh <- &StorageBucketSet{
		Bucket:   bIndex,
		ConnType: conntype, Addr: addr,
		Extra: extra,
	}
}

func (dw *distWorker) spin() {
	for {
		select {
		case msg := <-dw.msgCh:
			handleMsg(msg)
		case cmd := <-dw.clientCmdCh:
			mesh.SendAll(cmd)
		case sbs := <-dw.setBucketCh:
			mesh.SendAll(sbs)
		case <-dw.errCh:
			// TODO do something with the error
		}
	}
}

func handleMsg(msg interface{}) {
	switch m := msg.(type) {
	case StorageBucketSet:
		StorageBucketSetCommands <- &m
	case types.ClientCommand:
		ClientCommands <- &m
	}
}
