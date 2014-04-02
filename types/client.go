package types

// Action is the structure that commands from the client are parsed into.
// This is going to be the same regardless of the client interface or data
// format.
type Action struct {

	// Command gets passed through to the backend data-store.
	Command string `json:"cmd"`

	// StorageKey is the key used to route the command to the proper hyrax node.
	// Depending on the datastore backend it may also be incorporated into the
	// actual command sent to the datastore.
	StorageKey string `json:"key"`

	// Args are extra arguments needed for the command. This will depend on the
	// datastore used. The items in the args list can be of any type, but I
	// can't imagine needing anything except strings and numbers.
	Args []interface{} `json:"args,omitempty"`

	// Id is an optional identifier for who is sending this command.
	Id string `json:"id,omitempty"`

	// Secret is the sha1-hmac which is required for all commands which
	// add/change data in the datastore. The secret encompasses the command, the
	// key, and the id.
	Secret string `json:"secret,omitempty"`
}

// ActionReturn is the structure that returns to the client are parsed into.
// This is going to be the same regardless of the client interface or data
// format.
type ActionReturn struct {

	// Error will be filled out if there was an error somewhere in the command
	Error string `json:"error,omitempty"`

	// Return will be filled out if the command completed successfully. It will
	// be filled with whatever was returned from the command
	Return interface{} `json:"return,omitempty"`
}

// ErrorReturn takes in an error and returns a ActionReturn for it
func ErrorReturn(err error) *ActionReturn {
	return &ActionReturn{Error: err.Error()}
}
