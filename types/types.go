package types

import (
    "strconv"
)

type ConnId uint64
func (cid *ConnId) Serialize() string {
    return strconv.Itoa(int(*cid))
}

type Payload struct {
    Domain string   `json:"domain"`
    Id     string   `json:"id"`
    Name   string   `json:"name"`
    Secret string   `json:"secret"`
    Values []string `json:"values"`
}

type Command struct {
    Command string        `json:"command"`
    Payload Payload       `json:"payload"`
}
