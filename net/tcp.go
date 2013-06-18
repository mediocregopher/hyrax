package net

import (
    "net"
    "log"
    "bytes"
    "io"
    "hyrax/dispatch"
)

func TcpListen(addr string) error {
    log.Println("Starting tcp listener for",addr)
    listener, err := net.Listen("tcp",addr)
    if err != nil { return err }

    for i:=0; i<10; i++ {
        go func(){
            for {
                conn, err := listener.Accept()
                if err != nil {
                    log.Println("accept failed:",err.Error())
                    continue
                }

                cid,ch := AllocateId()
                go TcpClient(conn,cid,ch)
            }
        }()
    }

    return nil
}

type tcpReadChRet struct {
    msg []byte
    err error
}

func TcpClient(conn net.Conn, cid ConnId, cmdCh chan Message) {

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
                    //TODO does this ever happen?
                    ret = tcpReadChRet{nil,nil}
                }
                workerReadCh <- &ret
            }()
            readMore = false
        }

        select {

        //If we pull a command off we decode it and act accordingly
        case command := <-cmdCh:
            switch command.Type() {
            case PUSH:
                conn.Write(*command.(*PushMessage))
            }


        //If the goroutine doing the reading gets data we check it for an error
        //and send it to the globalReadCh to be handled
        case rcr := <-workerReadCh:
            readMore = true
            msg,err := rcr.msg,rcr.err
            if err != nil {
                TcpClose(conn,cid,cmdCh)
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
                        r,err := dispatch.DoCommand(fullMsg)
                        if err != nil {
                            log.Println("Go error from dispatch.DoCommand",err)
                            continue
                        }

                        conn.Write(r)
                        conn.Write([]byte{'\n'})
                    }
                }
            }

        }
    }
}

func TcpClose(conn net.Conn, cid ConnId, cmdCh chan Message) {
    conn.Close()
    CleanId(cid)
}
