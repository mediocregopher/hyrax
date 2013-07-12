package storage

import (
    "errors"
    "hyrax/config"
    "github.com/fzzy/radix/redis"
)

const NUM_CONNS = 10

type redisCmd struct {
    cmd *string
    args []interface{}
    ret chan *redis.Reply
}

var conns [NUM_CONNS]*redis.Client
var cmdCh chan *redisCmd

// init sets up NUM_CONNS routines that will each
// handle incoming commands for their channel
func init() {
    cmdCh = make(chan *redisCmd)
    for i := range conns {
        go func(i int){
            for cmd := range cmdCh {
                cmd.ret <- conns[i].Cmd( *cmd.cmd, cmd.args )
            }
        }(i)
    }
}

// RedisConnect creates all the connections for redis
func RedisConnect() error {
    addr := config.GetStr("redis-addr")

    for i := range conns {
        conn, err := redis.Dial("tcp",addr)
        if err != nil { return err }
        conns[i] = conn
    }
    return nil
}

// CmdPretty is a wrapper around Cmd, it takes in a command and a variadic
// list of arguments for that command. Used for commands defined directly
// in hyrax's code, where Cmd is used for commands that are being invoked
// by a client (generally)
func CmdPretty(cmd string, args... interface{}) (interface{},error) {
    return Cmd(cmd,args)
}

// Cmd takes in a command and a list of arguments for that command. It then
// determines the type of the return, decoding it, as well figuring out if there's
// been an error
func Cmd(cmd string, args []interface{}) (interface{},error) {
    rCmd := redisCmd{ &cmd, args, make(chan *redis.Reply) }
    cmdCh <- &rCmd
    r := <-rCmd.ret

    switch r.Type {
        case redis.StatusReply:
            return r.Str()

        case redis.ErrorReply:
            return nil,r.Err

        case redis.IntegerReply:
            return r.Int()

        case redis.NilReply:
            return nil,nil

        case redis.BulkReply:
            return r.Str()

        case redis.MultiReply:
            return r.List()
    }

    return nil,nil
}

// RedisListToMap is used for converting the return of a HGETALL to a hash
func RedisListToMap(l []string) (map[string]string,error) {
    llen := len(l)
    if llen%2 != 0 {
        return nil,errors.New("list has uneven number of elements")
    }

    m := map[string]string{}

    halfllen := llen/2
    for i := 0; i<halfllen; i++ {
        m[l[i*2]] = l[i*2+1]
    }

    return m,nil
}
