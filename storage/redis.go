package storage

import (
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
