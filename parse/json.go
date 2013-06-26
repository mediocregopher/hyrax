package parse

import (
    "encoding/json"
    "hyrax/types"
    "bytes"
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

// EncodeMessagePackage takes in a list of encoded messages as bytes
// and creates a proper json list
func EncodeMessagePackage(msgs [][]byte) ([]byte,error) {
    var buf bytes.Buffer
    buf.WriteByte('[')
    buf.Write(bytes.Join(msgs,[]byte{','}))
    buf.WriteByte(']')
    return buf.Bytes(),nil
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

// DecodeCommandPackage takes in raw json and tries to decode it into
// a list of commands
func DecodeCommandPackage(b []byte) ([]*types.Command,error) {
    var c []*types.Command
    err := json.Unmarshal(b,&c)
    return c,err
}
