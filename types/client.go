package types

// Storage key is used to identify a key that a command is taking place at (for
// example, SET, DEL, or any other redis command, assuming redis is your
// datastore).
type StorageKey []byte

// Implements the Bytes method for the Byter interface
func (s StorageKey) Bytes() []byte {
	return []byte(s)
}

// ClientCommand is the structure that commands from the client are parsed into.
// This is going to be the same regardless of the client interface or data
// format.
type ClientCommand struct {

	// Command gets passed through to the backend data-store.
	Command []byte `json:"cmd"`

	// StorageKey is the key used to route the command to the proper hyrax node.
	// Depending on the datastore backend it may also be incorporated into the
	// actual command sent to the datastore.
	StorageKey `json:"key"`

	// Args are extra arguments needed for the command. This will depend on the
	// datastore used. The items in the args list can be of any type, but I
	// can't imagine needing anything except strings and numbers.
	Args []interface{} `json:"args"`

	// Id is an optional identifier for who is sending this command.
	Id []byte `json:"id"`

	// Secret is the sha1-hmac which is required for all commands which
	// add/change data in the datastore. The secret encompasses the command, the
	// key, and the id.
	Secret []byte `json:"secret"`
}

// ClientReturn is the structure that returns to the client are parsed into.
// This is going to be the same regardless of the client interface or data
// format.
type ClientReturn struct {

	// Error will be filled out if there was an error somewhere in the command
	Error []byte `json:"error,omitempty"`

	// Return will be filled out if the command completed successfully. It will
	// be filled with whatever was returned from the command
	Return interface{} `json:"return,omitempty"`

}

// ErrorReturn takes in an error and returns a ClientReturn for it
func ErrorReturn(err error) *ClientReturn {
	return &ClientReturn{Error: []byte(err.Error())}
}
