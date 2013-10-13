package client

import (
	//See encoding/json/HACKED for why
	"github.com/mediocregopher/hyrax/encoding/json"
)

// ClientCommandToJson takes in a client command and returns the json form of
// it, or an error if something is whack.
func ClientCommandToJson(cc *ClientCommand) ([]byte, error) {
	return json.Marshal(cc)
}

// JsonToClientCommand takes in a byte slice, which is presumably filled with a
// json encoded ClientCommand, and decodes it.
func JsonToClientCommand(j []byte) (*ClientCommand, error) {
	cc := &ClientCommand{}
	err := json.Unmarshal(j, cc)
	return cc, err
}

// ClientReturnToJson takes in a client return struct and returns the json form
// of it, or an error if something is whack.
func ClientReturnToJson(cr *ClientReturn) ([]byte, error) {
	return json.Marshal(cr)
}

// JsonToClientReturn takes in a byte slice, which is presumably filled with a
// json encoded ClientReturn, and decodes it.
func JsonToClientReturn(j []byte) (*ClientReturn, error) {
	cr := &ClientReturn{}
	err := json.Unmarshal(j, cr)
	return cr, err
}
