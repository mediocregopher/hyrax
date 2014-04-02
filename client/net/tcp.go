package net

import (
	"bufio"
	"errors"
	"github.com/mediocregopher/manatcp"
	"io"

	"github.com/mediocregopher/hyrax/translate"
	"github.com/mediocregopher/hyrax/types"
)

type TcpClient struct {
	trans translate.Translator
	conn  *manatcp.Conn
}

func NewTcpClient(t translate.Translator, addr string,
	pushCh chan *types.Action) (*TcpClient, error) {

	tc := TcpClient{trans: t}
	conn, err := manatcp.Dial(&tc, addr)
	if err != nil {
		return nil, err
	}

	tc.conn = conn
	go func() {
		for ai := range conn.PushCh {
			if pushCh != nil {
				pushCh <- ai.(*types.Action)
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

	// Try to decode Action. We know it was a Action if Command is
	// actually set
	a, err := tc.trans.ToAction(b)
	if err != nil {
		return nil, err, false
	} else if a.Command != "" {
		return a, nil, false
	}

	ar, err := tc.trans.ToActionReturn(b)
	if err != nil {
		return nil, err, false
	}
	return ar, nil, false
}

func (tc *TcpClient) IsPush(ai interface{}) bool {
	_, ok := ai.(*types.Action)
	return ok
}

func (tc *TcpClient) Write(buf *bufio.Writer, ai interface{}) (error, bool) {
	b, err := tc.trans.FromAction(ai.(*types.Action))
	if err != nil {
		return err, false
	}

	if _, err := buf.Write(b); err != nil {
		return err, true
	}

	if _, err := buf.Write([]byte("\n")); err != nil {
		return err, true
	}

	return nil, false
}

func (tc *TcpClient) Cmd(cmd *types.Action) (interface{}, error) {
	ari, err, closed := tc.conn.Cmd(cmd)
	if closed {
		return nil, io.EOF
	} else if err != nil {
		return nil, err
	}

	ar, ok := ari.(*types.ActionReturn)
	if !ok {
		return nil, errors.New("Did not receive CommandReturn back")
	}

	if ar.Error != "" {
		return nil, errors.New(ar.Error)
	}

	return ar.Return, nil
}

func (tc *TcpClient) Close() {
	tc.conn.Close()
}
