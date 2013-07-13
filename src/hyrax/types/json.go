package types

import (
    "encoding/json"
    "bytes"
    "fmt"
)

// ByteSlice is a wrapper around a byte slice, which we
// use because the json marshaler wants to turn bytes into
// base64 strings
type ByteSlice []byte

func (b ByteSlice) MarshalJSON() ([]byte,error) {
    buf := make([]byte,0,len(b)+2)
    buf = append(buf,'"')
    buf = append(buf,b...)
    buf = append(buf,'"')
    return buf,nil
}

func (b *ByteSlice) UnmarshalJSON(json []byte) error {
    jlen := len(json)
    if json[0] != '"' || json[jlen-1] != '"' {
        return fmt.Errorf("%s is not a string",json)
    }
    *b = make([]byte,jlen-2)
    copy(*b,json[1:jlen-1])
    return nil
}


type Payload struct {
    Domain ByteSlice   `json:"domain"`
    Id     ByteSlice   `json:"id"`
    Name   ByteSlice   `json:"name"`
    Secret ByteSlice   `json:"secret"`
    Values []ByteSlice `json:"values"`
}


// Command (and subsequently Payload) are populated by json from the client and
// contain all relevant information about a command, so they're passed around a
// lot
type Command struct {
    Command ByteSlice  `json:"command"`
    Payload Payload    `json:"payload"`
    Quiet   bool       `json:"quiet"`
}

type messageWrap struct {
    Command ByteSlice     `json:"command"`
    Return  interface{}   `json:"return"`
}

type errorMessage struct {
    Command ByteSlice `json:"command,omitempty"`
    Error   ByteSlice `json:"error"`
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