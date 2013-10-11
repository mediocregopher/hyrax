package redis

import (
	"github.com/mediocregopher/hyrax/src/hyrax-server/router/storage/command"
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

