package net

import (
	"bufio"
	"io"
	"net"
	"errors"
	"github.com/mediocregopher/hyrax/types"
	"github.com/mediocregopher/hyrax/translate"
)

type cmdWrap struct {
	cmd *types.ClientCommand
	retCh chan *types.ClientReturn
}

type TcpClient struct {
	conn     net.Conn
	cmdCh    chan *cmdWrap
	closeCh  chan struct{}
	pushCh   chan *types.ClientCommand
	cmdRetCh chan *types.ClientReturn
	trans    translate.Translator
}

// Returns a new tcp client for hyrax which can be used to interact with a hyrax
// node. A push channel of nil indicates that you don't care about push messages
func NewTcpClient(
	t translate.Translator,
	addr string,
	pushCh chan *types.ClientCommand) (*TcpClient, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	c := &TcpClient{
		conn: conn,
		cmdCh: make(chan *cmdWrap),
		closeCh: make(chan struct{}),
		pushCh: pushCh,
		cmdRetCh: make(chan *types.ClientReturn),
		trans: t,
	}

	go c.readSpin()
	go c.writeSpin()
	return c, nil
}

// TODO I need to make a generic library for this crap
func (c *TcpClient) readSpin() {
	bufReader := bufio.NewReader(c.conn)

	for {
		b, err := bufReader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		cmd, err := c.trans.ToClientCommand(b)
		if err != nil {
			continue
		} else if cmd.Command != nil {
			if c.pushCh != nil {
				c.pushCh <- cmd
			}
			continue
		}

		ret, err := c.trans.ToClientReturn(b)
		if err != nil {
			c.cmdRetCh <- types.ErrorReturn(err)
		} else {
			c.cmdRetCh <- &types.ClientReturn{Return: ret}
		}
	}
	if c.pushCh != nil {
		close(c.pushCh)
	}
	close(c.cmdRetCh)
}

func (c *TcpClient) writeSpin() {
	spinloop: for {
		select {
		case <- c.closeCh:
			break spinloop
		case cw := <- c.cmdCh:
			b, err := c.trans.FromClientCommand(cw.cmd)
			if err != nil {
				cw.retCh <- types.ErrorReturn(err)
				break
			} else if err := c.write(b); err != nil {
				cw.retCh <- types.ErrorReturn(err)
				break
			}
			
			if ret, ok := <- c.cmdRetCh; ok {
				cw.retCh <- ret
			} else {
				cw.retCh <- types.ErrorReturn(io.EOF)
			}
		}
	}
	c.conn.Close()
}

func (c *TcpClient) write(b []byte) error {
	if _, err := c.conn.Write(b); err != nil {
		return err
	}

	if _, err := c.conn.Write([]byte("\n")); err != nil {
		return err
	}

	return nil
}

// Cmd sends the given command to the connection and returns the response
func (c *TcpClient) Cmd(cmd *types.ClientCommand) (interface{}, error) {
	cw := cmdWrap{cmd, make(chan *types.ClientReturn)}
	c.cmdCh <- &cw
	ret := <- cw.retCh
	if ret.Error != nil {
		return nil, errors.New(string(ret.Error))
	}
	return ret.Return, nil
}

// Close closes the tcp client
func (c *TcpClient) Close() {
	close(c.closeCh)
}
