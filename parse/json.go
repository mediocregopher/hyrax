package parse

import (
    "encoding/json"
    "hyrax/types"
)

type messageWrap struct {
    Command string        `json:"command"`
    Return  interface{}   `json:"return"`
}

type errorMessage struct {
    Command string `json:"command"`
    Error   string `json:"error"`
}

// EncodeMessage takes a return value from a given command,
// and returns the raw json
func EncodeMessage(command string, ret interface{}) ([]byte,error) {
    return json.Marshal(messageWrap{ command, ret })
}

// EncodeError takes in an error and returns the raw json for it
func EncodeError(command,err string) ([]byte,error) {
    return json.Marshal(errorMessage{command,err})
}

// DecodeCommand takes in raw json and tries to decode it into
// a command
func DecodeCommand(b []byte) (*types.Command,error) {
    var c types.Command
    err := json.Unmarshal(b,&c)
    return &c,err
}
