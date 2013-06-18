package net

import (
    "sync"
)

type MessageType int
const (
    PUSH MessageType = iota
)

type Message interface {
    Type() MessageType
}

type PushMessage []byte
func (m *PushMessage) Type() MessageType { return PUSH }

/////////////////////////////////////////////////////////

type ConnId uint64
type allocatedConn struct{ cid ConnId; ch chan Message }
var cidCh = make(chan *allocatedConn)

var router = struct{
    sync.RWMutex
    m map[ConnId]chan Message
}{
    m: map[ConnId]chan Message{},
}

func AllocateId() (ConnId, chan Message) {
    a := <-cidCh
    router.Lock()
    router.m[a.cid] = a.ch
    router.Unlock()
    return a.cid,a.ch
}

func CleanId(cid ConnId) {
    router.Lock()
    delete(router.m,cid)
    router.Unlock()
}

func GetChan(cid ConnId) (chan Message,bool) {
    router.RLock()
    ch,ok := router.m[cid]
    router.RUnlock()
    return ch,ok
}

func SendPushMessage(cid ConnId, msg []byte) bool {
    ch,ok := GetChan(cid)
    if !ok { return false }
    pmsg := PushMessage(msg)
    ch <- &pmsg
    return true
}

func init() {

    go func() {
        var i ConnId = 0
        for {
            cidCh <- &allocatedConn{ i, make(chan Message) }
            i++
        }
    }()

}
