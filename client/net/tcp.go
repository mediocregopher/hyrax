package net

import (
	"bufio"
	"io"
	"net"
	"errors"
	"github.com/mediocregopher/hyrax/types"
	"github.com/mediocregopher/hyrax/translate"
)

type TcpClient struct {
	conn     net.Conn
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
		pushCh: pushCh,
		cmdRetCh: make(chan *types.ClientReturn),
		trans: t,
	}

	go c.readSpin()
	return c, nil
}

// TODO I need to make a generic library for this crap
func (c *TcpClient) readSpin() {
	bufReader := bufio.NewReader(c.conn)

	for {
		b, err := bufReader.ReadBytes('\n')
		if err != nil {
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
			c.cmdRetCh <- ret
		}
	}
	c.conn.Close()
	if c.pushCh != nil {
		close(c.pushCh)
	}
	close(c.cmdRetCh)
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
	b, err := c.trans.FromClientCommand(cmd)
	if err != nil {
		return nil, err
	} else if err := c.write(b); err != nil {
		return nil, err
	}

	ret, ok := <- c.cmdRetCh
	if !ok {
		return nil, io.EOF
	} else if ret.Error != nil {
		return nil, errors.New(string(ret.Error))
	} else {
		return ret.Return, nil
	}
}

// Close closes the tcp client
func (c *TcpClient) Close() {
	c.conn.Close()
}
