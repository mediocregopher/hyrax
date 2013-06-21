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
    Error string `json:"error"`
}

func EncodeMessage(command string, ret interface{}) ([]byte,error) {
    return json.Marshal(messageWrap{ command, ret })
}

func EncodeError(err string) ([]byte,error) {
    return json.Marshal(errorMessage{err})
}

func DecodeCommand(b []byte) (*types.Command,error) {
    var c types.Command
    err := json.Unmarshal(b,&c)
    return &c,err
}
