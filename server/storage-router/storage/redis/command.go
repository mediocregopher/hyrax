package redis

import (
	"github.com/mediocregopher/hyrax/types"
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
var keyjoin = types.SimpleByter([]byte(":"))
var clientns = types.SimpleByter([]byte("cli"))

func (r *RedisKeyMaker) Namespace(ns, key types.Byter) types.Byter {
	return types.ByterJoin(keyjoin, ns, key)
}

func (r *RedisKeyMaker) ClientNamespace(ns, key types.Byter) types.Byter {
	return types.ByterJoin(keyjoin, clientns, ns, key)
}
