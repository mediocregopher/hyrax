package redis

import (
	"bytes"
	"github.com/mediocregopher/hyrax/server/storage-router/storage/command"
)

type RedisCommand struct {
	cmd  []byte
	args []interface{}
	trans []command.Command
}

func NewRedisCommand(cmd []byte, args []interface{}) command.Command {
	return &RedisCommand{
		cmd: cmd,
		args: args,
		trans: nil,
	}
}

func NewRedisTransaction(cmds ...command.Command) command.Command {
	return &RedisCommand{
		cmd: nil,
		args: nil,
		trans: cmds,
	}
}

func (c *RedisCommand) Cmd() []byte {
	return c.cmd
}

func (c *RedisCommand) Args() []interface{} {
	return c.args
}

func (c *RedisCommand) ExpandTransaction() []command.Command {
	return c.trans
}

type RedisKeyMaker struct {}
var keyjoin = []byte(":")
var clientns = []byte("cli")

func (r *RedisKeyMaker) Namespace(ns, key []byte) []byte {
	return bytes.Join([][]byte{ns, key}, keyjoin)
}

func (r *RedisKeyMaker) ClientNamespace(ns, key []byte) []byte {
	return bytes.Join([][]byte{clientns, ns, key}, keyjoin)
}
