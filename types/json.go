package types

import (
	"bytes"
	"hyrax/encoding/json"
)

// Command (and subsequently Payload) are populated by json from the client and
// contain all relevant information about a command, so they're passed around a
// lot
type Command struct {
	Command []byte  `json:"command"`
	Payload Payload `json:"payload"`
	Quiet   bool    `json:"quiet"`
}

type Payload struct {
	Domain []byte   `json:"domain"`
	Id     []byte   `json:"id"`
	Name   []byte   `json:"name"`
	Secret []byte   `json:"secret"`
	Values [][]byte `json:"values"`
}

type returnMessage struct {
	Command []byte      `json:"command"`
	Return  interface{} `json:"return"`
}

type errorMessage struct {
	Command []byte `json:"command,omitempty"`
	Error   []byte `json:"error"`
}

// EncodeReturn takes a return value from a given command, and returns the raw
// json
func EncodeReturn(command []byte, ret interface{}) ([]byte, error) {
	return json.Marshal(returnMessage{command, ret})
}

// EncodeReturnPackage takes in a list of encoded return messages as bytes and
// creates a proper json list
func EncodeReturnPackage(msgs [][]byte) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('[')
	buf.Write(bytes.Join(msgs, []byte{','}))
	buf.WriteByte(']')
	return buf.Bytes(), nil
}

// EncodeError takes in an error and returns the raw json for it
func EncodeError(command []byte, err error) ([]byte, error) {
	return json.Marshal(errorMessage{command, []byte(err.Error())})
}

// EncodeCommand takes in a command and returns its byte form
func EncodeCommand(cmd *Command) ([]byte,error) {
	return json.Marshal(cmd)
}

// DecodeCommand takes in raw json and tries to decode it into
// a command
func DecodeCommand(b []byte) (*Command, error) {
	var c Command
	err := json.Unmarshal(b, &c)
	return &c, err
}

// DecodeCommandPackage takes in raw json and tries to decode it into
// a list of commands
func DecodeCommandPackage(b []byte) ([]*Command, error) {
	var c []*Command
	err := json.Unmarshal(b, &c)
	return c, err
}
