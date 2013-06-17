package storage

import (
    "errors"
    "hyrax/config"
    "github.com/fzzy/radix/redis"
)

var conn *redis.Client

func RedisConnect() error {
    var err error
    addr := config.GetStr("redis-addr")
    conn, err = redis.Dial("tcp",addr)
    return err
}

func Cmd(cmd string, args ...interface{}) (interface{},error) {
    r := conn.Cmd(cmd,args...)
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
