package types

import (
    "strconv"
)

// ConnId is a unique value that's given to every connection on this hyrax node
type ConnId uint64
func (cid *ConnId) Serialize() string {
    return strconv.Itoa(int(*cid))
}

func ConnIdDeserialize(s string) (ConnId,error) {
    i,err := strconv.Atoi(s)
    if err != nil { return 0,err }
    return ConnId(i),nil
}

// MonPushPayload is the payload for push notifications. It is basically
// the standard payload object but without the secret, and with a command
// string field instead
type MonPushPayload struct {
    Domain  string   `json:"domain"`
    Id      string   `json:"id"`
    Name    string   `json:"name,omitempty"`
    Command string   `json:"command"`
    Values  []string `json:"values,omitempty"`
}

