package command

import (
	"github.com/mediocregopher/hyrax/types"
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

	// If this command represents a transaction then ExpandTransaction returns
	// all the sub-commands that the transaction encompasses. If this command
	// isn't a transaction then this returns nil
	ExpandTransaction() []Command
}

// CommandBundle is a grouping of a Command and a channel on which the
// CommandRet for that command will be returned.
type CommandBundle struct {
	Cmd Command
	RetCh chan *CommandRet
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

	// DirectCommandAllowed returns whether or not a direct command is allowed
	// to be executed by a client
	DirectCommandAllowed(cmd []byte) bool

	// DirectCommandModifies returns whether or not a direct command modifies
	// existing state in the storage unit
	DirectCommandModifies(cmd []byte) bool

	// KeyValueSets are sets of (innerkey -> value) mappings located at key.
	// They are queried, added, and removed by innerkey. It's also possible to
	// get a list of all values being held. KeyValueSetAdd is a command which
	// adds an (innerkey -> value) map (or overwrites it, if innerky was already
	// mapped to something). It creates the set at key with this one mapping if
	// it didn't exist previously.
	KeyValueSetAdd(
		key types.Byter,
		innerKey types.Byter,
		value types.Byter) Command

	// KeyValueSetRemByInnerKey removes an (innerkey -> value) mapping from the
	// keyval set located at key, or does nothing if that key, or that innerkey
	// within the key, didn't exist.
	KeyValueSetRemByInnerKey(key types.Byter, innerkey types.Byter) Command

	// KeyValueSetCard is a command which returns the number of (innerkey ->
	// value) mappings are in the set at the given key. It should return a zero
	// value if the given key doesn't exist
	KeyValueSetCard(key types.Byter) Command

	// KeyValueSetMemberValues is a command which returns the list of value
	// portions of the (innerkey -> value) mappings in the set. It should return
	// an empty list if the set doesn't exist
	KeyValueSetMemberValues(key types.Byter) Command

	// KeyValueSetDel is a command which deletes the entire set at the given
	// key.  Nothing should happen if the set doesn't exist in the first place
	KeyValueSetDel(key types.Byter) Command

	// A generic set is your run-of-the-mill set, where each value in the set
	// only appears once. GenericSetAdd adds a value to the set at key, and
	// creates the set if it didn't previously exist.
	GenericSetAdd(key, value types.Byter) Command

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

	// GenericSetDel is a command which deletes the entire set at the given key.
	// Nothing should happen if the set doesn't exist in the first place
	GenericSetDel(key types.Byter) Command
}

// KeyMaker is responsible taking in keys for clients and creating new
// namespaced keys for them. Depending on the datastore this may not be needed,
// but for ones like redis it is.
type KeyMaker interface{

	// Namespace takes a namespace and a key and returns a combination of them
	// which makes sense for the storage backend. For some backends, this may
	// simply ignore the namespace
	Namespace(ns, key types.Byter) types.Byter

	// ConnNamespace is the same as Namespace except that it is used for
	// metadata related to a specific client, so it may be formatted a bit
	// different then Namespace's output for the same inputs
	ClientNamespace(ns, key types.Byter) types.Byter

}
