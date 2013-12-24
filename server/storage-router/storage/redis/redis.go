package redis

import (
	"fmt"
	"github.com/fzzy/radix/redis"
	"github.com/mediocregopher/hyrax/server/storage-router/storage/command"
	"github.com/mediocregopher/hyrax/server/storage-router/storage/unit"
	"log"
	"time"
)

type RedisConn struct {
	conn    *redis.Client
	cmdCh   chan *command.CommandBundle
	closeCh chan chan error
}

func New() unit.StorageUnitConn {
	return &RedisConn{
		conn:    nil,
		cmdCh:   make(chan *command.CommandBundle),
		closeCh: make(chan chan error),
	}
}

// Implements Connect for StorageUnitConn. Connects to redis over tcp and spawns
// a handler go-routine
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
			ret := command.CommandRet{rawret, err}
			select {
			case cmdb.RetCh <- &ret:
			case <-time.After(10 * time.Second):
				log.Printf("Timedout in redisconn replying to cmd %v", cmdb.Cmd)
			}
		}
	}
}

func (r *RedisConn) cmd(cmd command.Command) (interface{}, error) {
	if trans := cmd.ExpandTransaction(); trans != nil {
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
		reply := r.conn.Cmd(string(cmd.Cmd()), cmd.Args()...)
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
func (r *RedisConn) Cmd(cmdb *command.CommandBundle) {
	select {
	case r.cmdCh <- cmdb:
	case <-time.After(10 * time.Second):
		err := fmt.Errorf("Redis connection timedout receiving command")
		select {
		case cmdb.RetCh <- &command.CommandRet{nil, err}:
		case <-time.After(1 * time.Second):
		}
	}
}

// Implements Close for StorageUnitConn.
func (r *RedisConn) Close() error {
	retCh := make(chan error)
	r.closeCh <- retCh
	return <-retCh
}
