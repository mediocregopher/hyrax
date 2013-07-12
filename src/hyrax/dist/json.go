package dist

import (
    "encoding/json"
    "hyrax/types"
    "errors"
)

type messageType int
const (
    MONPUSH messageType = iota
)

type messageIn struct {
    Type    messageType
    Payload json.RawMessage
}

type messageOut struct {
    Type    messageType
    Payload interface{}
}

// DecodeMessage takes in raw json and tries to decode it into a message
func DecodeMessage(b []byte) (interface{},error) {
    var m messageIn
    err := json.Unmarshal(b,&m)
    if err != nil { return nil,err}

    switch m.Type {
        case MONPUSH:
            var r types.MonPushPayload
            err = json.Unmarshal(m.Payload,&r)
            return &r,err
    }

    return nil,errors.New("Unknown message type in dist.DecodeMessage")
}

// EncodeMessage takes in a messageOut object and encodes it into bytes
func EncodeMessage(msg *messageOut) ([]byte,error) {
    return json.Marshal(msg)
}
