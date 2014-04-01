package redis

import (
	"fmt"
	"github.com/fzzy/radix/redis"
	"github.com/grooveshark/golib/gslog"
	"time"

	"github.com/mediocregopher/hyrax/server/storage"
)

// A connection to redis, implements Storage interface
type RedisConn struct {
	conn    *redis.Client
	cmdCh   chan *storage.CommandBundle
	closeCh chan chan error
}

// A command into the redis connection, implements the Command interface
type RedisCommand struct {
	cmd  string
	args []interface{}
}

// Returns a new RedisCommand with the given arguments
func NewRedisCommand(cmd string, args []interface{}) storage.Command {
	return &RedisCommand{
		cmd:  cmd,
		args: args,
	}
}

func (c *RedisCommand) Cmd() string {
	return c.cmd
}

func (c *RedisCommand) Args() []interface{} {
	return c.args
}

// Returns an unconnected redis connection structure as per the Storage
// interface
func New() storage.Storage {
	return &RedisConn{
		conn:    nil,
		cmdCh:   make(chan *storage.CommandBundle),
		closeCh: make(chan chan error),
	}
}

// Implements Connect for Storage. Connects to redis over tcp and spawns a
// handler go-routine
func (r *RedisConn) Connect(conntype, addr string, _ ...interface{}) error {
	conn, err := redis.Dial(conntype, addr)
	if err != nil {
		return err
	}

	r.conn = conn
	go r.spin()
	return nil
}

func (r *RedisConn) spin() {
	for {
		select {

		case retCh := <-r.closeCh:
			retCh <- r.conn.Close()
			close(r.cmdCh)
			close(r.closeCh)
			return

		case cmdb := <-r.cmdCh:
			rawret, err := r.cmd(cmdb.Cmd)
			ret := storage.CommandRet{rawret, err}
			select {
			case cmdb.RetCh <- &ret:
			case <-time.After(10 * time.Second):
				gslog.Errorf("RedisConn timedout replying to cmd %v", cmdb.Cmd)
			}
		}
	}
}

func (r *RedisConn) cmd(cmd storage.Command) (interface{}, error) {
	reply := r.conn.Cmd(cmd.Cmd(), cmd.Args()...)
	return decodeReply(reply)
}

// Decodes a reply into a generic interface object, or an error
func decodeReply(r *redis.Reply) (interface{}, error) {
	switch r.Type {
	case redis.StatusReply:
		return r.Bytes()

	case redis.ErrorReply:
		return nil, r.Err

	case redis.IntegerReply:
		return r.Int()

	case redis.NilReply:
		return nil, nil

	case redis.BulkReply:
		return r.Bytes()

	case redis.MultiReply:
		return r.ListBytes()
	}

	return nil, nil
}

// Implements Cmd for Storage
func (r *RedisConn) Cmd(cmdb *storage.CommandBundle) {
	select {
	case r.cmdCh <- cmdb:
	case <-time.After(10 * time.Second):
		err := fmt.Errorf("Redis connection timedout receiving command")
		select {
		case cmdb.RetCh <- &storage.CommandRet{nil, err}:
		case <-time.After(1 * time.Second):
		}
	}
}

// Implements NewCommand for Storage
func (_ *RedisConn) NewCommand(cmd string, args []interface{}) storage.Command {
	return NewRedisCommand(cmd, args)
}

// Implements CommandAllowed for Storage
func (_ *RedisConn) CommandAllowed(cmd string) bool {
	_, ok := getCommandInfo(cmd)
	return ok
}

// Implements CommandModifies for Storage
func (_ *RedisConn) CommandModifies(cmd string) bool {
	cinfo, ok := getCommandInfo(cmd)
	return ok && cinfo.Modifies
}

// Implements CommandIsAdmin for Storage. Redis has no administrative commands
// which are allowed so this is always false
func (_ *RedisConn) CommandIsAdmin(_ string) bool {
	return false
}

// Implements Close for Storage
func (r *RedisConn) Close() error {
	retCh := make(chan error)
	r.closeCh <- retCh
	return <-retCh
}
