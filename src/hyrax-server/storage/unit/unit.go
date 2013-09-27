package unit

import (
	"errors"
	sucmd "github.com/mediocregopher/hyrax/src/hyrax-server/storage/command"
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
	// returning the result on the ret channel. This method shouldn't worry
	// about timeouts, the StorageUnit will worry about that.
	Cmd(cmd *sucmd.Command, ret chan *sucmd.CommandRet)

	// Close tells the connection that it's no longer needed. It should close
	// any external resources it has open and tell all internal go-routines to
	// end execution. Any subsequent calls to the StorageUnitConn will cause a
	// panic.
	Close() error
}

type storageUnitCmd struct {
	cmd *sucmd.Command
	ret chan *sucmd.CommandRet
}

// A storage unit is a pool of storage unit conns which can be opened and closed
// as a single group. It also multiplexes calls across the connections.
type StorageUnit struct {
	conns []StorageUnitConn
	cmdCh chan *storageUnitCmd
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
		cmdCh: make(chan *storageUnitCmd),
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

			case sucmd := <-su.cmdCh:
				// This is not great. It could lead to a race-condition if the
				// call to Cmd (which is sending on a channel to the connection
				// go-routine presumably) comes in AFTER a call to close on that
				// same connection, which could happen in the next iteration of
				// this loop. I'm not sure of any good solutions to this other
				// then to make this happen synchronously, which would be really
				// bad for performance if one connection suddenly locks up.
				go su.conns[i].Cmd(sucmd.cmd, sucmd.ret)
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

	cmdS := sucmd.Command{
		Cmd: cmd,
		Args: args,
	}

	sucmd := storageUnitCmd{
		cmd: &cmdS,
		ret: make(chan *sucmd.CommandRet),
	}
	select {
	case su.cmdCh <- &sucmd:
	case <- time.After(10 * time.Second):
		return nil, errors.New("sending command to StorageUnit timedout")
	}

	select {
	case ret := <- sucmd.ret:
		return ret.Ret, ret.Err
	case <- time.After(10 * time.Second):
		return nil, errors.New("receiving results from StorageUnit timedout")
	}

}
