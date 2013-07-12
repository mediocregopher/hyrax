package dist

import (
    "net"
    "log"
    "bytes"
    "io"
    "hyrax/types"
    "hyrax/custom"
)

// DistListen starts up a tcp socket which will listen for dist messages
func DistListen(addr string) error {
    log.Println("Starting dist listener for",addr)
    listener, err := net.Listen("tcp",addr)
    if err != nil { return err }

    go func(){
        for {
            conn, err := listener.Accept()
            if err != nil {
                log.Println("dist accept failed:",err.Error())
                continue
            }

            go DistClient(conn)
        }
    }()

    return nil
}

type tcpReadChRet struct {
    msg []byte
    err error
}

// DistClient handles the actual work of reading commands from a connection
func DistClient(conn net.Conn) {

    workerReadCh  := make(chan *tcpReadChRet)
    readMore := true
    readBuf := new(bytes.Buffer)
    for {

        //readMore keeps track of whether or not a routine is already reading
        //off the connection. If there isn't one we make another
        if readMore {
            go func(){
                var ret tcpReadChRet
                buf := make([]byte,1024)
                bcount, err := conn.Read(buf)
                if err != nil {
                    ret = tcpReadChRet{nil,err}
                } else if bcount > 0 {
                    ret = tcpReadChRet{buf,nil}
                } else {
                    ret = tcpReadChRet{nil,nil}
                }
                workerReadCh <- &ret
            }()
            readMore = false
        }

        select {

        //If the goroutine doing the reading gets data we check it for an error
        //and send it to the globalReadCh to be handled
        case rcr := <-workerReadCh:
            readMore = true
            msg,err := rcr.msg,rcr.err
            if err != nil {
                log.Println("dist connection with "+conn.RemoteAddr().String()+" died")
                conn.Close()
                return
            } else if msg != nil {
                readBuf.Write(msg)
                for {
                    fullMsg,err := readBuf.ReadBytes('\n')
                    if err == io.EOF {
                        //We got to the end of the buffer without finding a delim,
                        //write back what we did find so it can be searched the next time
                        readBuf.Reset()
                        if fullMsg[0] != '\x00' {
                            readBuf.Write(fullMsg)
                        }
                        break
                    } else {
                        go DistDispatch(fullMsg)
                    }
                }
            }
        }
    }
}

// DistDispatch takes in bytes that are presumably a message,
// decodes them, and does whatever they say to do
func DistDispatch(b []byte) {
    r,err := DecodeMessage(b)
    if err != nil {
        log.Println(err.Error()+" in dist.DistDispatch when decoding")
        return
    }

    switch msg := r.(type) {
        case *types.MonPushPayload:
            err = custom.MonDoAlert(msg)
            if err != nil {
                log.Println(err.Error()+" in dist.DistDispatch when calling custom.MonDoAlert")
            }
        default:
            log.Println("Unknown message type in dist.DistDispatch")
    }
}
