package command

import (
	"github.com/mediocregopher/hyrax/src/hyrax/types"
)


// CommandRet is returned from a Command in the RetCh. It's really just a tuple
// around the return value and an error
type CommandRet struct {
	Ret interface{}
	Err error
}

// Command is sent to a StorageUnitConn, and contains all data necessary to
// complete a call and return any data from it. Both methods implemented aren't
// really used outside of the CommandFactory and StorageUnitConn they're being
// created for, so their behavior isn't terribly important.
type Command interface {

	// Cmd returns the basic command that's going to be executed. This will
	// differ depending on backend
	Cmd() []byte

	// Args is a list of arguments to the command. This will also differ
	// depending on platform.
	Args() []interface{}
}

// A CommandFactory (it's not a real factory in the OO sense, I just couldn't
// think of a better name) is a set of methods that need to be implemented by a
// backend data-store in order for the rest of hyrax-server to correctly
// interface with it. Each method creates a command to perform that specific
// task, and returns the fully initialized Command struct for that command.
type CommandFactory interface {

	// Transaction takes in multiple commands that have been returned from this
	// factory and returns a command which represents an atomic transaction of
	// all of them in the order provided. The way this works will probably need
	// to change for compatibility with other backed storages besides redis
	Transaction(...Command) Command

	// DirectCommand is for commands that the client is calling directly on the
	// storage for itself.
	DirectCommand(cmd []byte, key types.Byter, args []interface{}) Command

	// IdValueSets are sets of (id,value) tuples, which can be queried for both
	// by id and value. IdValueSetAdd is a command which adds a tuple (id,value)
	// to the set at key, and creates that set if it didn't exist before.
	IdValueSetAdd(
		key types.Byter,
		id types.Uint64er,
		value types.Byter) Command

	// IdValueSetRem removes an (id,value) tuple from an IdValueSet. There is no
	// error if the set did not exist.
	IdValueSetRem(
		key types.Byter,
		id types.Uint64er,
		value types.Byter) Command

	// IdValueSetIsIdMember is a command which returns whether or not a given id
	// is a member in the set at key. The command created should return a falsey
	// value if the set didn't exist.
	IdValueSetIsIdMember(key types.Byter, id types.Uint64er) Command

	// IdValueSetIsValueMember is a command which returns whether or not a given
	// value is a member in the set at key. The command created should return a
	// falsey value if the set didn't exist.
	IdValueSetIsValueMember(key, value types.Byter) Command

	// IdValueSetCard is a command which returns the number of (id,value) tuples
	// are in the set at the given key. It should return a zero value if the
	// given key doesn't exist
	IdValueSetCard(key types.Byter) Command

	// IdValueSetMemberValues is a command which returns the list of value
	// portions of the (id,value) tuples in the set. It should return an empty
	// list if the set doesn't exist
	IdValueSetMemberValues(key types.Byter) Command

	// A generic set is your run-of-the-mill set, where each value in the set
	// only appears once. GenericSetAdd adds a value to the set at key, and
	// creates the set if it didn't previously exist.
	GenericSetAdd(key, value types.Byter)

	// GenericSetRem removes the given value from the set at the given key.
	// There is no error if the set did not exist.
	GenericSetRem(key, value types.Byter) Command

	// GenericSetIsMember is a command which returns whether or not the given
	// value is in the set at the given key. The command should return a falsey
	// value if the set did not exist.
	GenericSetIsMember(key, value types.Byter) Command

	// GenericSetCard is a command which returns the number of values in the set
	// at the given key. It should return a zero value if the set doesn't exist.
	GenericSetCard(key types.Byter) Command

	// GenericSetMembers is a command which returns the list of values in the
	// set at the given key. It should return an empty list if the set doesn't
	// exist.
	GenericSetMembers(key types.Byter) Command
}
