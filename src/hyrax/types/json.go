package types

import (
    "encoding/json"
    "bytes"
)

type Payload struct {
    Domain string   `json:"domain"`
    Id     string   `json:"id"`
    Name   string   `json:"name"`
    Secret string   `json:"secret"`
    Values []string `json:"values"`
}


// Command (and subsequently Payload) are populated by json from the client and
// contain all relevant information about a command, so they're passed around a
// lot
type Command struct {
    Command string  `json:"command"`
    Payload Payload `json:"payload"`
    Quiet   bool    `json:"quiet"`
}

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
func DecodeCommand(b []byte) (*Command,error) {
    var c Command
    err := json.Unmarshal(b,&c)
    return &c,err
}

// DecodeCommandPackage takes in raw json and tries to decode it into
// a list of commands
func DecodeCommandPackage(b []byte) ([]*Command,error) {
    var c []*Command
    err := json.Unmarshal(b,&c)
    return c,err
}
