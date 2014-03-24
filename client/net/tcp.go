package net

import (
	"bufio"
	"errors"
	"io"
	"github.com/mediocregopher/manatcp"

	"github.com/mediocregopher/hyrax/translate"
	"github.com/mediocregopher/hyrax/types"
)

type TcpClient struct {
	trans translate.Translator
	conn  *manatcp.Conn
}

func NewTcpClient(t translate.Translator, addr string,
	pushCh chan *types.ClientCommand) (*TcpClient, error) {
	
	tc := TcpClient{trans: t}
	conn, err := manatcp.Dial(&tc, addr)
	if err != nil {
		return nil, err
	}

	tc.conn = conn
	go func() {
		for cci := range conn.PushCh {
			if pushCh != nil {
				pushCh <- cci.(*types.ClientCommand)
			}
		}
	}()
	return &tc, nil
}

func (tc *TcpClient) Read(buf *bufio.Reader) (interface{}, error, bool) {
	b, err := buf.ReadBytes('\n')
	if err != nil {
		return nil, err, true
	}

	// Try to decode ClientCommand. We know it was a ClientCommand if Command is
	// actually set
	cc, err := tc.trans.ToClientCommand(b)
	if err != nil {
		return nil, err, false
	} else if cc.Command != nil {
		return cc, nil, false
	}

	cr, err := tc.trans.ToClientReturn(b)
	if err != nil {
		return nil, err, false
	}
	return cr, nil, false
}

func (tc *TcpClient) IsPush(lc interface{}) bool {
	_, ok := lc.(*types.ClientCommand)
	return ok
}

func (tc *TcpClient) Write(buf *bufio.Writer, item interface{}) (error, bool) {
	b, err := tc.trans.FromClientCommand(item.(*types.ClientCommand))
	if err != nil {
		return err, false
	}

	if _, err := buf.Write(b); err != nil {
		return err, true
	}

	return nil, false
}

func (tc *TcpClient) Cmd(cmd *types.ClientCommand) (interface{}, error) {
	cri, err, closed := tc.conn.Cmd(cmd)
	if closed {
		return nil, io.EOF
	} else if err != nil {
		return nil, err
	}

	cr, ok := cri.(*types.ClientReturn)
	if !ok {
		return nil, errors.New("Did not receive CommandReturn back")
	}

	if cr.Error != nil {
		return nil, errors.New(string(cr.Error))
	}

	return cr.Return, nil
}

func (tc *TcpClient) Close() {
	tc.conn.Close()
}
