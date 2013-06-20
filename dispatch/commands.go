package dispatch

import (
    "hyrax/types"
    "hyrax/custom"
)

type commandInfo uint
const (
    CUSTOM commandInfo = 1 << iota
    MODIFY
    RETURNS_MAP
)

func CommandExists(cmd *string) bool {
    _,ok := commandMap[*cmd]
    return ok
}

func CommandIsCustom(cmd *string) bool {
    n := commandMap[*cmd]
    return n & CUSTOM > 0
}

func CommandModifies(cmd *string) bool {
    n := commandMap[*cmd]
    return n & MODIFY > 0
}

func CommandReturnsMap(cmd *string) bool {
    n := commandMap[*cmd]
    return n & RETURNS_MAP > 0
}

var commandMap = map[string]commandInfo{

    //Keys
    "exists":           0,
    "expire":           MODIFY,
    "expireat":         MODIFY,
    "persist":          MODIFY,
    "pexpire":          MODIFY,
    "pexpireat":        MODIFY,
    "pttl":             0,
    "ttl":              0,
    "type":             0,

    //Strings
    "append":           MODIFY,
    "bitcount":         0,
    "decr":             MODIFY,
    "decrby":           MODIFY,
    "get":              0,
    "getbit":           0,
    "getrange":         0,
    "getset":           MODIFY,
    "incr":             MODIFY,
    "incrby":           MODIFY,
    "incrbyfloat":      MODIFY,
    "psetex":           MODIFY,
    "set":              MODIFY,
    "setbit":           MODIFY,
    "setex":            MODIFY,
    "setnx":            MODIFY,
    "setrange":         MODIFY,
    "strlen":           0,

    //Hashes
    "hdel":             MODIFY,
    "hexists":          0,
    "hget":             0,
    "hgetall":          RETURNS_MAP,
    "hincrby":          MODIFY,
    "hincrbyfloat":     MODIFY,
    "hkeys":            0,
    "hlen":             0,
    "hmget":            0,
    "hset":             MODIFY,
    "hsetnx":           MODIFY,
    "hvals":            0,

    //Lists
    //blpop
    //brpop
    "lindex":           0,
    "linsert":          MODIFY,
    "llen":             0,
    "lpop":             MODIFY,
    "lpush":            MODIFY,
    "lpushx":           MODIFY,
    "lrange":           0,
    "lrem":             MODIFY,
    "lset":             MODIFY,
    "ltrim":            MODIFY,
    "rpop":             MODIFY,
    "rpush":            MODIFY,
    "rpushx":           MODIFY,

    //Sets
    "sadd":             MODIFY,
    "scard":            0,
    "sismember":        0,
    "smembers":         0,
    "spop":             MODIFY,
    "srandmember":      0,
    "srem":             MODIFY,

    //Sorted Sets
    "zadd":             MODIFY,
    "zcard":            0,
    "zcount":           0,
    "zincrby":          MODIFY,
    "zrange":           0,
    "zrangebyscore":    0,
    "zrank":            0, //TODO two different return values?
    "zrem":             MODIFY,
    "zremrangebyrank":  MODIFY,
    "zremrangebyscore": MODIFY,
    "zrevrange":        0,
    "zrevrangebyscore": 0,
    "zrevrank":         0, //TODO two different return values?
    "zscore":           0,

    //Monitors
    "mon":              CUSTOM,
    "hmon":             CUSTOM,
    "lmon":             CUSTOM,
    "smon":             CUSTOM,
    "zmon":             CUSTOM,
    "amon":             CUSTOM,
    "emon":             CUSTOM,
}

var customCommandMap = map[string]func(types.ConnId, *types.Payload)(interface{},error){
    "amon":             custom.AMon,
    "mon":              custom.Mon,
    "hmon":             custom.HMon,
    "lmon":             custom.LMon,
    "smon":             custom.SMon,
    "zmon":             custom.ZMon,
}
