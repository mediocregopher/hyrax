package types

import (
    "strconv"
)

// ConnId is a unique value that's given to every connection on this hyrax node
type ConnId uint64
func (cid *ConnId) Serialize() string {
    return strconv.Itoa(int(*cid))
}

// Command (and subsequently Payload) are populated by json from the client and
// contain all relevant information about a command, so they're passed around a
// lot

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
