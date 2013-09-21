package storageunit

import (
	"errors"
	"time"
)

// A storage unit connection is a connection to a single storage unit, which
// acts like a basic tcp connection but could in theory be just about anything.
type StorageUnitConn interface {

	// Connect is called on a zero'd StorageUnitConn and causes it to create its
	// initial connection to the datastore and set up any internal go-routines
	// that are needed.
	Connect(conntype, addr string, extra ...interface{}) error

	// Cmd takes in a command struct and processes the command contained within,
	// as well as responding on the RetCh in it. This method shouldn't worry
	// about timeouts, the StorageUnit will worry about that.
	Cmd(cmd *Command)

	// Close tells the connection that it's no longer needed. It should close
	// any external resources it has open and tell all internal go-routines to
	// end execution. Any subsequent calls to the StorageUnitConn will cause a
	// panic.
	Close() error
}

// CommandRet is returned from a Command in the RetCh. It's really just a tuple
// around the return value and an error
type CommandRet struct {
	Ret interface{}
	Err error
}

// Command is sent to a StorageUnitConn, and contains all data necessary to
// complete a call and return any data from it.
type Command struct {
	Cmd   []byte
	Args  []interface{}
	RetCh chan *CommandRet
}

// A storage unit is a pool of storage unit conns which can be opened and closed
// as a single group. It also multiplexes calls across the connections.
type StorageUnit struct {
	conns []StorageUnitConn
	cmdCh chan *Command
	closeCh chan chan error
}

// NewStorageUnit takes in a slice of zero'd StorageUnitConns, a connection
// type/address, and any extra info which will be passed back to the Connect
// function of each StorageUnitConn. If any of the calls to Connect result in an
// error all the previous StorageUnitConns will be Close'd and that error will
// be returned.
func NewStorageUnit(
	sucs []StorageUnitConn, conntype, addr string, extra ...interface{}) error {
	
	su := StorageUnit{
		conns: make([]StorageUnitConn, 0, len(sucs)),
		cmdCh: make(chan *Command),
		closeCh: make(chan chan error),
	}

	for _, suc := range sucs {
		if err := suc.Connect(conntype, addr, extra...); err == nil {
			su.conns = append(su.conns, suc)
		} else {
			su.internalClose()
			return err
		}
	}

	go su.spin()
	return nil
}

func (su *StorageUnit) spin() {
	for {
		for i := range su.conns {
			select {

			case retCh := <-su.closeCh:
				retCh <- su.internalClose()
				close(su.closeCh)
				close(su.cmdCh)
				return

			case cmd := <-su.cmdCh:
				su.conns[i].Cmd(cmd)
			}
		}
	}
}

func (su *StorageUnit) internalClose() error {
	var retErr error
	for i := range su.conns {
		if err := su.conns[i].Close(); err != nil {
			retErr = err
		}
	}
	return retErr
}

// Close goes through all StorageUnitConns that are currently active and calls
// Close() on each one. The last non-nil error to be returned by any of those
// Close calls is returned, or nil if none of them returned an error. Any
// subsequent calls to the StorageUnit will cause a panic.
func (su *StorageUnit) Close() error {
	retCh := make(chan error)
	su.closeCh <- retCh
	return <- retCh
}

// Cmd takes in a command and its arguments, performs the command in one of the
// StorageUnitConn's, and returns any return data from performing the command.
func (su *StorageUnit) Cmd(
	cmd []byte, args []interface{}) (interface{}, error) {

	cmdS := Command{
		Cmd: cmd,
		Args: args,
		RetCh: make(chan *CommandRet),
	}

	select {
	case su.cmdCh <- &cmdS:
	case <- time.After(10 * time.Second):
		return nil, errors.New("sending command to StorageUnit timedout")
	}

	select {
	case ret := <- cmdS.RetCh:
		return ret.Ret, ret.Err
	case <- time.After(10 * time.Second):
		return nil, errors.New("receiving results from StorageUnit timedout")
	}

}
