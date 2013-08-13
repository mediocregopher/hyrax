package storage

import (
	"errors"
	"github.com/fzzy/radix/redis"
	"hyrax-server/config"
)

const NUM_CONNS = 10

type redisCmd struct {
	cmd  *string
	args []interface{}
	ret  chan *redis.Reply
}

var conns [NUM_CONNS]*redis.Client
var cmdCh chan *redisCmd

// init sets up NUM_CONNS routines that will each
// handle incoming commands for their channel
func init() {
	cmdCh = make(chan *redisCmd)
	for i := range conns {
		go func(i int) {
			for cmd := range cmdCh {
				cmd.ret <- conns[i].Cmd(*cmd.cmd, cmd.args)
			}
		}(i)
	}
}

// RedisConnect creates all the connections for redis
func RedisConnect() error {
	addr := config.GetStr("redis-addr")

	for i := range conns {
		conn, err := redis.Dial("tcp", addr)
		if err != nil {
			return err
		}
		conns[i] = conn
	}
	return nil
}

// CmdPretty is a wrapper around Cmd, it takes in a command and a variadic
// list of arguments for that command. Used for commands defined directly
// in hyrax's code, where Cmd is used for commands that are being invoked
// by a client (generally)
func CmdPretty(cmd []byte, args ...interface{}) (interface{}, error) {
	return Cmd(cmd, args)
}

// Cmd takes in a command and a list of arguments for that command. It then
// determines the type of the return, decoding it, as well figuring out if there's
// been an error
func Cmd(cmd []byte, args []interface{}) (interface{}, error) {
	scmd := string(cmd)
	rCmd := redisCmd{&scmd, args, make(chan *redis.Reply)}
	cmdCh <- &rCmd
	r := <-rCmd.ret

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
		//I have a merge request with radix for a proper
		//ListBytes() method. Hopefully the guy notices,
		//if not I'll change the radix dependency to be
		//my fork
		l, err := r.List()
		if err != nil {
			return nil, err
		}
		lb := make([][]byte, len(l))
		for i := range l {
			if l[i] == "" {
				lb[i] = nil
			} else {
				lb[i] = []byte(l[i])
			}
		}
		return lb, nil
	}

	return nil, nil
}

// RedisListToMap is used for converting the return of a HGETALL to a hash
func RedisListToMap(l [][]byte) (map[string][]byte, error) {
	llen := len(l)
	if llen%2 != 0 {
		return nil, errors.New("list has uneven number of elements")
	}

	m := map[string][]byte{}

	halfllen := llen / 2
	for i := 0; i < halfllen; i++ {
		m[string(l[i*2])] = l[i*2+1]
	}

	return m, nil
}
