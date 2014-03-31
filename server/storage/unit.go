package storage

import (
	"fmt"
	"time"
)

// CommandRet is returned from a Command in the RetCh. It's really just a tuple
// around the return value and an error
type CommandRet struct {
	Ret interface{}
	Err error
}

// Command is sent to a Storage, and contains all data necessary to
// complete a call and return any data from it. Both methods implemented aren't
// really used outside of the CommandFactory and Storage they're being
// created for, so their behavior isn't terribly important.
type Command interface {

	// Cmd returns the basic command that's going to be executed. This will
	// differ depending on backend
	Cmd() []byte

	// Args is a list of arguments to the command. This will also differ
	// depending on platform.
	Args() []interface{}
}

// CommandBundle is a grouping of a Command and a channel on which the
// CommandRet for that command will be returned.
type CommandBundle struct {
	Cmd   Command
	RetCh chan *CommandRet
}

// A storage unit connection is a connection to a single storage unit, which
// acts like a basic tcp connection but could in theory be just about anything.
type Storage interface {

	// Connect is called on a zero'd Storage and causes it to create its
	// initial connection to the datastore and set up any internal go-routines
	// that are needed.
	Connect(conntype, addr string, extra ...interface{}) error

	// Cmd takes in a command bundle and processes the command contained within,
	// returning the result on the ret channel. If the implementation of this
	// method involves any blocking operations then they should have a timeout
	// which, when hit, stops the command and sends an error on the ret channel
	Cmd(*CommandBundle)

	// Given the command and arguments for an action on the datastore, returns a
	// Command instance. This method should not actually affect anything about
	// the Storage connection.
	NewCommand([]byte, []interface{}) Command

	// Returns whether or not a command is allowed to be called under any
	// circumstances. This method should not actually affect anything about the
	// Storage connection.
	CommandAllowed([]byte) bool

	// Returns whether or not a command will modify state within the datastore
	// (and therefore potentially require authentication). This method should
	// not actually affect anything about the Storage connection.
	CommandModifies([]byte) bool

	// Returns whether or not a command requires administrative privileges (and
	// therefore potentially require authentication). This method should not
	// actually affect anything about the Storage connection.
	CommandIsAdmin([]byte) bool

	// Close tells the connection that it's no longer needed. It should close
	// any external resources it has open and tell all internal go-routines to
	// end execution. Any subsequent calls to the Storage will cause a
	// panic.
	Close() error
}

// A storage unit is a pool of storage unit conns which can be opened and closed
// as a single group. It also multiplexes calls across the connections.
type StorageUnit struct {
	ConnType, Addr string
	conns          []Storage
	cmdCh          chan *CommandBundle
	closeCh        chan chan error
}

// NewStorageUnit takes in a slice of zero'd Storages, a connection
// type/address, and any extra info which will be passed back to the Connect
// function of each Storage. If any of the calls to Connect result in an
// error all the previous Storages will be Close'd and that error will
// be returned.
func NewStorageUnit(
	sucs []Storage,
	conntype, addr string,
	extra ...interface{}) (*StorageUnit, error) {

	su := StorageUnit{
		ConnType: conntype,
		Addr:     addr,
		conns:    make([]Storage, 0, len(sucs)),
		cmdCh:    make(chan *CommandBundle),
		closeCh:  make(chan chan error),
	}

	for _, suc := range sucs {
		if err := suc.Connect(conntype, addr, extra...); err == nil {
			su.conns = append(su.conns, suc)
		} else {
			su.internalClose()
			return nil, err
		}
	}

	go su.spin()
	return &su, nil
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
				// TODO FIX THIS
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

// Close goes through all Storages that are currently active and calls
// Close() on each one. The last non-nil error to be returned by any of those
// Close calls is returned, or nil if none of them returned an error. Any
// subsequent calls to the StorageUnit will cause a panic.
func (su *StorageUnit) Close() error {
	retCh := make(chan error)
	su.closeCh <- retCh
	return <-retCh
}

// Cmd takes in a command bundle and performs the command in one of the
// Storage's connections
func (su *StorageUnit) Cmd(cmd Command) (interface{}, error) {
	cmdb := &CommandBundle{cmd, make(chan *CommandRet)}
	select {
	case su.cmdCh <- cmdb:
	case <-time.After(10 * time.Second):
		// TODO logging
		err := fmt.Errorf(
			"sending command to StorageUnit %s:%s timedout",
			su.ConnType,
			su.Addr,
		)
		return nil, err
	}

	select {
	case cret := <-cmdb.RetCh:
		return cret.Ret, cret.Err
	case <-time.After(10 * time.Second):
		// TODO logging
		err := fmt.Errorf(
			"receiving response from StorageUnit %s:%s timedout",
			su.ConnType,
			su.Addr,
		)
		return nil, err
	}
}

// Returns a new Command instance based on the given command and arguments
func (su *StorageUnit) NewCommand(cmd []byte, args []interface{}) Command {
	return su.conns[0].NewCommand(cmd, args)
}

// Returns whether or not a command is allowed to be run at all on the datastore
func (su *StorageUnit) CommandAllowed(cmd []byte) bool {
	return su.conns[0].CommandAllowed(cmd)
}

// Returns whether or not a command modifies state on the datastore
func (su *StorageUnit) CommandModifies(cmd []byte) bool {
	return su.conns[0].CommandModifies(cmd)
}

func (su *StorageUnit) CommandIsAdmin(cmd []byte) bool {
	return su.conns[0].CommandIsAdmin(cmd)
}
