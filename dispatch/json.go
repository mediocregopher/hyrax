package dispatch

import (
    "encoding/json"
)

type CommandPart struct {
    Domain string   `json:"domain"`
    Id     string   `json:"id"`
    Name   string   `json:"name"`
    Secret string   `json:"secret"`
    Values []string `json:"values"`
}

type Command struct {
    Command string        `json:"command"`
    Payload []CommandPart `json:"payload"`
}

//messageWrap and errorMessage are messages
type message interface{}

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

func DecodeCommand(b []byte) (*Command,error) {
    var c Command
    err := json.Unmarshal(b,&c)
    return &c,err
}
