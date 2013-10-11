package unit

import (
	"fmt"
	"github.com/mediocregopher/hyrax/src/hyrax-server/router/storage/command"
	"time"
)

// A storage unit connection is a connection to a single storage unit, which
// acts like a basic tcp connection but could in theory be just about anything.
type StorageUnitConn interface {

	// Connect is called on a zero'd StorageUnitConn and causes it to create its
	// initial connection to the datastore and set up any internal go-routines
	// that are needed.
	Connect(conntype, addr string, extra ...interface{}) error

	// Cmd takes in a command bundle and processes the command contained within,
	// returning the result on the ret channel. If the implementation of this
	// method involves any blocking operations then they should have a timeout
	// which, when hit, stops the command and sends an error on the ret channel
	Cmd(*command.CommandBundle)

	// Close tells the connection that it's no longer needed. It should close
	// any external resources it has open and tell all internal go-routines to
	// end execution. Any subsequent calls to the StorageUnitConn will cause a
	// panic.
	Close() error
}

// A storage unit is a pool of storage unit conns which can be opened and closed
// as a single group. It also multiplexes calls across the connections.
type StorageUnit struct {
	ConnType, Addr string
	conns []StorageUnitConn
	cmdCh chan *command.CommandBundle
	closeCh chan chan error
}

// NewStorageUnit takes in a slice of zero'd StorageUnitConns, a connection
// type/address, and any extra info which will be passed back to the Connect
// function of each StorageUnitConn. If any of the calls to Connect result in an
// error all the previous StorageUnitConns will be Close'd and that error will
// be returned.
func NewStorageUnit(
	sucs []StorageUnitConn,
	conntype, addr string,
	extra ...interface{}) (*StorageUnit,error) {
	
	su := StorageUnit{
		ConnType: conntype,
		Addr: addr,
		conns: make([]StorageUnitConn, 0, len(sucs)),
		cmdCh: make(chan *command.CommandBundle),
		closeCh: make(chan chan error),
	}

	for _, suc := range sucs {
		if err := suc.Connect(conntype, addr, extra...); err == nil {
			su.conns = append(su.conns, suc)
		} else {
			su.internalClose()
			return nil,err
		}
	}

	go su.spin()
	return &su,nil
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

			case command := <-su.cmdCh:
				// This is not great. It could lead to a race-condition if the
				// call to Cmd (which is sending on a channel to the connection
				// go-routine presumably) comes in AFTER a call to close on that
				// same connection, which could happen in the next iteration of
				// this loop. I'm not sure of any good solutions to this other
				// then to make this happen synchronously, which would be really
				// bad for performance if one connection suddenly locks up.
				go su.conns[i].Cmd(command)
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

// Cmd takes in a command bundle and performs the command in one of the
// StorageUnitConn's
func (su *StorageUnit) Cmd(cmdb *command.CommandBundle) {
	select {
	case su.cmdCh <- cmdb:
	case <- time.After(10 * time.Second):
		err := fmt.Errorf(
			"sending command to StorageUnit %s:%s timedout",
			su.ConnType,
			su.Addr,
		)
		select {
		case cmdb.RetCh <- &command.CommandRet{nil, err}:
		case <-time.After(1 * time.Second):
		}
	}
}
