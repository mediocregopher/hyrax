package router

import (
    "sync"
    . "hyrax-server/types" //Import all cause we're gonna need ConnId alot!
)

// MessageType's are the different types of messages that
// a connection routine will handle
type MessageType int
const (
    PUSH MessageType = iota
)

// Message's are pushed to connection routines telling them to do stuff
type Message interface {
    Type() MessageType
}

// PushMessage is data that the connection routine should forward to the
// connection
type PushMessage []byte
func (m *PushMessage) Type() MessageType { return PUSH }

/////////////////////////////////////////////////////////

type allocatedConn struct{ cid ConnId; ch chan Message }
var cidCh = make(chan *allocatedConn)

// router keeps track of the connection id to connection Message channel
// mappings
var router = struct{
    sync.RWMutex
    m map[ConnId]chan Message
}{
    m: map[ConnId]chan Message{},
}

// AllocateId creates a new id and message channel for use by a connection
func AllocateId() (ConnId, chan Message) {
    a := <-cidCh
    router.Lock()
    router.m[a.cid] = a.ch
    router.Unlock()
    return a.cid,a.ch
}

// CleanId removes any lingering data from a dead connection in this module
func CleanId(cid ConnId) {
    router.Lock()
    delete(router.m,cid)
    router.Unlock()
}

// getChan returns the message channel for a given connection id
func getChan(cid ConnId) (chan Message,bool) {
    router.RLock()
    ch,ok := router.m[cid]
    router.RUnlock()
    return ch,ok
}

// SendPushMessage sends a message to a connection routine telling it to
// send the given data to a connection
func SendPushMessage(cid ConnId, msg []byte) bool {
    ch,ok := getChan(cid)
    if !ok { return false }
    pmsg := PushMessage(msg)
    ch <- &pmsg
    return true
}

// init sets up a routine that dishes out connection ids like a mofo!
func init() {
    go func() {
        var i ConnId = 0
        for {
            cidCh <- &allocatedConn{ i, make(chan Message) }
            i++
        }
    }()
}
