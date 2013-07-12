package types

import (
    "encoding/json"
    "bytes"
)

// byteSlice is a wrapper around a byte slice, which we
// use because the json marshaler wants to turn bytes into
// base64 strings
type byteSlice []byte
func (b byteSlice) MarshalJSON() ([]byte,error) {
    buf := make([]byte,0,len(b)+2)
    buf = append(buf,'"')
    buf = append(buf,b...)
    buf = append(buf,'"')
    return buf,nil
}


type Payload struct {
    Domain byteSlice   `json:"domain"`
    Id     byteSlice   `json:"id"`
    Name   byteSlice   `json:"name"`
    Secret byteSlice   `json:"secret"`
    Values []byteSlice `json:"values"`
}


// Command (and subsequently Payload) are populated by json from the client and
// contain all relevant information about a command, so they're passed around a
// lot
type Command struct {
    Command byteSlice  `json:"command"`
    Payload Payload    `json:"payload"`
    Quiet   bool       `json:"quiet"`
}

type messageWrap struct {
    Command byteSlice        `json:"command"`
    Return  interface{}   `json:"return"`
}

type errorMessage struct {
    Command byteSlice `json:"command"`
    Error   byteSlice `json:"error"`
}

// EncodeMessage takes a return value from a given command,
// and returns the raw json
func EncodeMessage(command []byte, ret interface{}) ([]byte,error) {
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
func EncodeError(command,err []byte) ([]byte,error) {
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
