package dist

import (
    "net"
    "sync"
    "log"
    "errors"
)

var conns = map[string]chan []byte{}
var connLock = sync.RWMutex{}

// Connect makes a dist connection to addr
func Connect(addr string) error {
    conn,err := net.Dial("tcp",addr)
    if err != nil { return err }

    go putter(conn,addr)
    return nil
}

// Disconnect disables the dist connection to addr
func Disconnect(addr string) {
    connLock.RLock()
    ch,ok := conns[addr]
    connLock.RUnlock()
    if !ok { return }

    connLock.Lock()
    delete(conns,addr)
    connLock.Unlock()
    close(ch)
}

func putter(conn net.Conn, addr string) {
    ch := make(chan []byte)
    connLock.Lock()
    conns[addr] = ch
    connLock.Unlock()

    var err error
    for msg := range ch {

        _,err = conn.Write(msg)
        if err != nil {
            log.Println(err.Error()+" in dist.putter when calling conn.Write")
            connLock.Lock()
            delete(conns,addr)
            connLock.Unlock()
            return
        }

        _,err = conn.Write([]byte{'\n'})
        if err != nil {
            log.Println(err.Error()+" in dist.putter when calling conn.Write")
            connLock.Lock()
            delete(conns,addr)
            connLock.Unlock()
            return
        }
    }
}

// Send sends a message to a specific node
func Send(msg []byte, addr string) error {
    connLock.RLock()
    ch,ok := conns[addr]
    connLock.RUnlock()
    if !ok {
        return errors.New("Address "+addr+" not connected in dist")
    }

    ch <- msg
    return nil
}

// Broadcast sends a message to all connected nodes
func Broadcast(msg []byte) {
    connLock.RLock()
    for _,ch := range conns {
        ch <- msg
    }
    connLock.RUnlock()
}
