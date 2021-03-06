package redis

import (
	"github.com/fzzy/radix/redis"
	"github.com/grooveshark/golib/gslog"
	"io"
	"time"

	"github.com/mediocregopher/hyrax/server/storage"
)

// A command into the redis connection, implements the Command interface
type RedisCommand struct {
	cmd  string
	args []interface{}
}

// Returns a new RedisCommand with the given arguments
func NewRedisCommand(cmd string, args ...interface{}) storage.Command {
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

////////////////////////////////////////////////////////////////////////////////

// A connection to redis, implements Storage interface
type RedisConn struct {
	conntype string
	addr     string
	conn     *redis.Client
	cmdCh    chan *storage.CommandBundle
	closeCh  chan chan error
}

// Returns an unconnected redis connection structure as per the Storage
// interface
func New() storage.Storage {
	return &RedisConn{}
}

// Implements Connect for Storage. Connects to redis over tcp and spawns a
// handler go-routine
func (r *RedisConn) Connect(cmdCh chan *storage.CommandBundle,
	conntype, addr string, _ ...interface{}) error {
	conn, err := redis.Dial(conntype, addr)
	if err != nil {
		gslog.Errorf("connecting to redis at %s: %s", addr, err)
		return err
	}

	r.conntype = conntype
	r.addr = addr
	r.conn = conn
	r.cmdCh = cmdCh
	r.closeCh = make(chan chan error)
	go r.spin()
	return nil
}

func (r *RedisConn) spin() {
spinloop:
	for {
		select {

		case retCh := <-r.closeCh:
			retCh <- r.conn.Close()
			break spinloop

		case cmdb := <-r.cmdCh:
			rawret, err := r.cmd(cmdb.Cmd)
			ret := storage.CommandRet{rawret, err}
			select {
			case cmdb.RetCh <- &ret:
			case <-time.After(10 * time.Second):
				gslog.Errorf("RedisConn timedout replying to cmd %v", cmdb.Cmd)
			}
			if err == io.EOF && !r.resurrect() {
				break spinloop
			}
		}
	}

	close(r.closeCh)
}

func (r *RedisConn) resurrect() bool {
	connCh := make(chan *redis.Client)
	go func() {
		for {
			time.Sleep(2 * time.Second)
			conn, err := redis.Dial(r.conntype, r.addr)
			if err != nil {
				gslog.Errorf("connecting to redis at %s: %s", r.addr, err)
				continue
			}
			connCh <- conn
			return
		}
	}()

	for {
		select {
		case retCh := <-r.closeCh:
			retCh <- r.conn.Close()
			return false
		case conn := <-connCh:
			r.conn = conn
			return true
		}
	}
}

func (r *RedisConn) cmd(cmd storage.Command) (interface{}, error) {
	gslog.Debugf("Redis cmd: %v, %v", cmd.Cmd(), cmd.Args())
	reply := r.conn.Cmd(cmd.Cmd(), cmd.Args()...)
	dreply, err := decodeReply(reply)
	gslog.Debugf("Redis reply: %v, %v", dreply, err)
	return dreply, err
}

// Decodes a reply into a generic interface object, or an error
func decodeReply(r *redis.Reply) (interface{}, error) {
	switch r.Type {
	case redis.StatusReply:
		gslog.Debugf("Redis status reply")
		return r.Str()

	case redis.ErrorReply:
		gslog.Debugf("Redis error reply")
		return nil, r.Err

	case redis.IntegerReply:
		gslog.Debugf("Redis int reply")
		return r.Int()

	case redis.NilReply:
		gslog.Debugf("Redis nil reply")
		return nil, nil

	case redis.BulkReply:
		gslog.Debugf("Redis bulk reply")
		return r.Str()

	case redis.MultiReply:
		gslog.Debugf("Redis multibulk reply")
		return r.List()
	}

	return nil, nil
}

// Implements NewCommand for Storage
func (_ *RedisConn) NewCommand(cmd string, args ...interface{}) storage.Command {
	return NewRedisCommand(cmd, args...)
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
