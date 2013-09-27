package redis

import (
	"errors"
	"github.com/fzzy/radix/redis"
	sucmd "github.com/mediocregopher/hyrax/src/hyrax-server/storage/command"
)

type RedisConn struct {
	conn *redis.Client
	cmdCh chan *cmdWrap
	closeCh chan chan error
}

type cmdWrap struct {
	cmd sucmd.Command
	ret chan *sucmd.CommandRet
}

func New() *RedisConn {
	return &RedisConn{
		conn: nil,
		cmdCh: make(chan *cmdWrap),
		closeCh: make(chan chan error),
	}
}

// Implements Connect for StorageUnitConn. Connects to redis over tcp and spawns
// a handler go-routine
func (r *RedisConn) Connect(conntype, addr string, _ ... interface{}) error {
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

		case retCh := <- r.closeCh:
			retCh <- r.conn.Close()
			close(r.cmdCh)
			close(r.closeCh)
			return

		case cmd := <- r.cmdCh:
			r, err := r.cmd(cmd)
			ret := sucmd.CommandRet{r, err}
			cmd.ret <- &ret
		}
	}
}

func (r *RedisConn) cmd(cmdwrap *cmdWrap) (interface{}, error) {
	if trans := cmdwrap.cmd.ExpandTransaction(); trans != nil {
		r.conn.Append(string(MULTI))
		for i := range trans {
			r.conn.Append(string(trans[i].Cmd()), trans[i].Args()...)
		}
		r.conn.Append(string(EXEC))

		for i := 0; i < len(trans)+1; i++ {
			r.conn.GetReply()
		}
		return decodeReply(r.conn.GetReply())
	} else {
		reply := r.conn.Cmd(string(cmdwrap.cmd.Cmd()), cmdwrap.cmd.Args()...)
		return decodeReply(reply)
	}
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

// Implements Cmd for StorageUnitConn.
func (r *RedisConn) Cmd(cmd sucmd.Command, cmdret chan *sucmd.CommandRet) {
	r.cmdCh <- &cmdWrap{cmd, cmdret}
}

// Implements Close for StorageUnitConn.
func (r *RedisConn) Close() error {
	retCh := make(chan error)
	r.closeCh <- retCh
	return <- retCh
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