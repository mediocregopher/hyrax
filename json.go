package main

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

const (
    STRING = iota
    INT
    FLOAT
    LIST
    MAP
)

type ReturnString []string
type ReturnInt    []int64
type ReturnFloat  []float64
type ReturnList   [][]string
type ReturnMap    []map[string]string

type MessageWrap struct {
    Command string        `json:"command"`
    Return  interface{}   `json:"return"`
}

type ErrorMessage struct {
    Error string `json:"error"`
}

func EncodeMessage(command string, ret interface{}) ([]byte,error) {
    return json.Marshal(MessageWrap{ command, ret })
}

func EncodeError(err string) ([]byte,error) {
    return json.Marshal(ErrorMessage{err})
}

func DecodeCommand(b []byte) (*Command,error) {
    var c Command
    err := json.Unmarshal(b,&c)
    return &c,err
}
