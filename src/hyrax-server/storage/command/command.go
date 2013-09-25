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
// complete a call and return any data from it.
type Command struct {
	Cmd   []byte
	Args  []interface{}
}

// A CommandFactory (it's not a real factory in the OO sense, I just couldn't
// think of a better name) is a set of methods that need to be implemented by a
// backend data-store in order for the reset of hyrax-server to correctly
// interface with it.
type CommandFactory interface {

	// DirectCommand is for commands that the client is calling directly on the
	// storage for itself.
	DirectCommand(cmd []byte, key types.Byter, args []interface{})

	// IdValueSets are sets of (id,value) tuples, which can be queried for both
	// by id and value. IdValueSetAdd is a command which adds a tuple (id,value)
	// to the set at key, and creates that set if it didn't exist before.
	IdValueSetAdd(key types.Byter, id types.Uint64er, value types.Byter)

	// IdValueSetRem removes an (id,value) tuple from an IdValueSet. There is no
	// error if the set did not exist.
	IdValueSetRem(key types.Byter, id types.Uint64er, value types.Byter)

	// IdValueSetIsIdMember is a command which returns whether or not a given id
	// is a member in the set at key. The command created should return a falsey
	// value if the set didn't exist.
	IdValueSetIsIdMember(key types.Byter, id types.Uint64er)

	// IdValueSetIsValueMember is a command which returns whether or not a given
	// value is a member in the set at key. The command created should return a
	// falsey value if the set didn't exist.
	IdValueSetIsValueMember(key, value types.Byter)

	// IdValueSetCard is a command which returns the number of (id,value) tuples
	// are in the set at the given key. It should return a zero value if the
	// given key doesn't exist
	IdValueSetCard(key types.Byter)

	// A generic set is your run-of-the-mill set, where each value in the set
	// only appears once. GenericSetAdd adds a value to the set at key, and
	// creates the set if it didn't previously exist.
	GenericSetAdd(key, value types.Byter)

	// GenericSetRem removes the given value from the set at the given key.
	// There is no error if the set did not exist.
	GenericSetRem(key, value types.Byter)

	// GenericSetIsMember is a command which returns whether or not the given
	// value is in the set at the given key. The command should return a falsey
	// value if the set did not exist.
	GenericSetIsMember(key, value types.Byter)

	// GenericSetCard is a command which return the number of values in the set
	// at the given key. It should return a zero value if the set doesn't exist.
	GenericSetCard(key types.Byter)
}
