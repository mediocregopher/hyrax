package dispatch

import (
    "hyrax/types"
    "hyrax/custom"
)

type commandInfo uint
const (
    CUSTOM commandInfo = 1 << iota
    MODIFY
    QUIET
    RETURNS_MAP
)

// CommandExists returns back whether or not
// a command is actually avaible for clients to use
func CommandExists(cmd *string) bool {
    _,ok := commandMap[*cmd]
    return ok
}

// CommandIsCustom returns back whether or not
// a command is a "custom" command, i.e: if the
// command is implemented in hyrax and isn't passed right
// through to redis.
func CommandIsCustom(cmd *string) bool {
    n := commandMap[*cmd]
    return n & CUSTOM > 0
}

// CommandModifies returns back whether or not a command
// modifies and existing value.
func CommandModifies(cmd *string) bool {
    n := commandMap[*cmd]
    return n & MODIFY > 0
}

// CommandIsQuiet returns back if a command which modifies
// a value should remain quiet and not generate an alert
func CommandIsQuiet(cmd *string) bool {
    n := commandMap[*cmd]
    return n & QUIET > 0
}

// CommandReturnsMap returns back whether or not a command
// is expected to return back a string->string map from redis
func CommandReturnsMap(cmd *string) bool {
    n := commandMap[*cmd]
    return n & RETURNS_MAP > 0
}

// commandMap is a map of commands to their bitmasks, where
// the bitmask describes the attributes of the command.
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
    "zrank":            0,
    "zrem":             MODIFY,
    "zremrangebyrank":  MODIFY,
    "zremrangebyscore": MODIFY,
    "zrevrange":        0,
    "zrevrangebyscore": 0,
    "zrevrank":         0,
    "zscore":           0,

    //Monitors
    "mon":              CUSTOM,
    "hmon":             CUSTOM,
    "lmon":             CUSTOM,
    "smon":             CUSTOM,
    "zmon":             CUSTOM,
    "amon":             CUSTOM,
    "emon":             CUSTOM,

    //EKGs
    "eadd":             MODIFY | CUSTOM,
    "eaddq":            MODIFY | CUSTOM | QUIET,
    "erem":             MODIFY | CUSTOM,
    "eremq":            MODIFY | CUSTOM | QUIET,
    "emembers":         CUSTOM,
}

// customCommandMap is a map of custom commands to their appropriate
// built-in funcions.
var customCommandMap = map[string]func(types.ConnId, *types.Payload)(interface{},error){

    //Monitors
    "amon":             custom.AMon,
    "mon":              custom.Mon,
    "hmon":             custom.HMon,
    "lmon":             custom.LMon,
    "smon":             custom.SMon,
    "zmon":             custom.ZMon,
    "emon":             custom.EMon,

    //EKGs
    "eadd":             custom.EAdd,
    "eaddq":            custom.EAdd,
    "erem":             custom.ERem,
    "eremq":            custom.ERem,
    "emembers":         custom.EMembers,

}
