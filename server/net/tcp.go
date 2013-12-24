package net

import (
	"bufio"
	"github.com/mediocregopher/hyrax/server/client"
	crouter "github.com/mediocregopher/hyrax/server/client-router"
	stypes "github.com/mediocregopher/hyrax/server/types"
	"github.com/mediocregopher/hyrax/translate"
	"github.com/mediocregopher/hyrax/types"
	"net"
)

// TcpClient holds a tcp connection into the hyrax server, and implements the
// Client interface
type TcpClient struct {
	id     stypes.ClientId
	pushCh chan *types.ClientCommand
	conn   net.Conn
	trans  translate.Translator
}

// Returns a new TcpClient structure, having populated it with all necessary
// fields and added it to the client router
func NewTcpClient(conn net.Conn, t translate.Translator) (*TcpClient, error) {
	cid := client.NewClient()
	c := &TcpClient{
		id:     cid,
		pushCh: make(chan *types.ClientCommand),
		conn:   conn,
		trans:  t,
	}
	if err := crouter.Add(c); err != nil {
		return nil, err
	}

	return c, nil
}

func (tc *TcpClient) ClientId() stypes.ClientId {
	return tc.id
}

func (tc *TcpClient) CommandPushChannel() chan<- *types.ClientCommand {
	return tc.pushCh
}

func (tc *TcpClient) write(b []byte) error {
	//TODO write deadline
	if _, err := tc.conn.Write(b); err != nil {
		return err
	}

	if _, err := tc.conn.Write([]byte("\n")); err != nil {
		return err
	}

	return nil
}

func (tc *TcpClient) writeClientReturn(cr *types.ClientReturn) error {
	b, err := tc.trans.FromClientReturn(cr)
	if err != nil {
		return err
	}

	return tc.write(b)
}

func (tc *TcpClient) writeErr(err error) error {
	cr := types.ClientReturn{Error: []byte(err.Error())}
	return tc.writeClientReturn(&cr)
}

type tcpReadChRet struct {
	msg []byte
	err error
}

func (tc *TcpClient) Spin() {
	workerReadCh := make(chan *tcpReadChRet)
	readMore := true
	bufReader := bufio.NewReader(tc.conn)
spin:
	for {

		//readMore keeps track of whether or not a routine is already reading
		//off the connection. If there isn't one we make another
		if readMore {
			go func() {
				b, err := bufReader.ReadBytes('\n')
				workerReadCh <- &tcpReadChRet{b, err}
			}()
			readMore = false
		}

		select {

		case cmd := <-tc.pushCh:
			b, _ := tc.trans.FromClientCommand(cmd)
			if err := tc.write(b); err != nil {
				break spin
			}

		case rcr := <-workerReadCh:
			readMore = true
			msg, err := rcr.msg, rcr.err
			if err != nil {
				break spin
			} else if msg != nil {
				cc, err := tc.trans.ToClientCommand(msg)
				if err != nil {
					if err = tc.writeErr(err); err != nil {
						break spin
					}
					continue
				}

				cr := client.RunCommand(tc.id, cc)
				if err = tc.writeClientReturn(cr); err != nil {
					break spin
				}
			}
		}
	}

	tc.conn.Close()
	client.ClientClosed(tc.id)
}

func TcpListen(addr string, trans translate.Translator) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			// TODO log error
			continue
		}

		c, err := NewTcpClient(conn, trans)
		if err != nil {
			// TODO log error
			continue
		}

		go c.Spin()
	}
}
