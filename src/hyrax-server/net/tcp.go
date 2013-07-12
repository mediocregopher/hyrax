package net

import (
    "net"
    "log"
    "bufio"
    "hyrax-server/types"
    "hyrax-server/dispatch"
    "hyrax-server/router"
)

// TcpListen starts up a tcp listen server, and sets up the acceptor routines
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

                cid,ch := router.AllocateId()
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

// TcpClient is the main function tcp connections use. It constantly reads
// in data from the network and the messsage push channel
func TcpClient(conn net.Conn, cid types.ConnId, cmdCh chan router.Message) {

    workerReadCh  := make(chan *tcpReadChRet)
    readMore := true
    bufReader := bufio.NewReader(conn)
    for {

        //readMore keeps track of whether or not a routine is already reading
        //off the connection. If there isn't one we make another
        if readMore {
            go func(){
                b,err := bufReader.ReadBytes('\n')
                workerReadCh <- &tcpReadChRet{b,err}
            }()
            readMore = false
        }

        select {

        //If we pull a command off we decode it and act accordingly
        case command := <-cmdCh:
            switch command.Type() {
            case router.PUSH:
                conn.Write(*command.(*router.PushMessage))
                conn.Write([]byte{'\n'})
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
                r,err := dispatch.DoCommand(cid,msg)
                if err != nil {
                    log.Println("Go error from dispatch.DoCommand:",err)
                    continue
                }

                conn.Write(r)
                conn.Write([]byte{'\n'})
            }
        }
    }
}

// TcpClose is used when a connection needs to be closed or when it's already been closed.
// It initiates cleanup of the connection and its data.
func TcpClose(conn net.Conn, cid types.ConnId, cmdCh chan router.Message) {
    conn.Close()
    err := dispatch.DoCleanup(cid)
    if err != nil { log.Println("Error during cleanup of",cid,":",err.Error()) }
}
