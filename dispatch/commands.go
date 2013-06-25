package dispatch

import (
    "hyrax/types"
    "hyrax/custom"
)

// CommandInfo is a struct which is tied to a command,
// and describes various of the command. All properties
// are false be default.
type CommandInfo struct {
    IsCustom,
    Modifies,
    IsQuiet,
    ReturnsMap bool
}

// commandMap is a map of commands to their info structs
var commandMap = map[string]*CommandInfo{

    //Keys
    "exists":           &CommandInfo{},
    "expire":           &CommandInfo{Modifies:true},
    "expireat":         &CommandInfo{Modifies:true},
    "persist":          &CommandInfo{Modifies:true},
    "pexpire":          &CommandInfo{Modifies:true},
    "pexpireat":        &CommandInfo{Modifies:true},
    "pttl":             &CommandInfo{},
    "ttl":              &CommandInfo{},
    "type":             &CommandInfo{},

    //Strings
    "append":           &CommandInfo{Modifies:true},
    "bitcount":         &CommandInfo{},
    "decr":             &CommandInfo{Modifies:true},
    "decrby":           &CommandInfo{Modifies:true},
    "get":              &CommandInfo{},
    "getbit":           &CommandInfo{},
    "getrange":         &CommandInfo{},
    "getset":           &CommandInfo{Modifies:true},
    "incr":             &CommandInfo{Modifies:true},
    "incrby":           &CommandInfo{Modifies:true},
    "incrbyfloat":      &CommandInfo{Modifies:true},
    "psetex":           &CommandInfo{Modifies:true},
    "set":              &CommandInfo{Modifies:true},
    "setbit":           &CommandInfo{Modifies:true},
    "setex":            &CommandInfo{Modifies:true},
    "setnx":            &CommandInfo{Modifies:true},
    "setrange":         &CommandInfo{Modifies:true},
    "strlen":           &CommandInfo{},

    //Hashes
    "hdel":             &CommandInfo{Modifies:true},
    "hexists":          &CommandInfo{},
    "hget":             &CommandInfo{},
    "hgetall":          &CommandInfo{ReturnsMap:true},
    "hincrby":          &CommandInfo{Modifies:true},
    "hincrbyfloat":     &CommandInfo{Modifies:true},
    "hkeys":            &CommandInfo{},
    "hlen":             &CommandInfo{},
    "hmget":            &CommandInfo{},
    "hset":             &CommandInfo{Modifies:true},
    "hsetnx":           &CommandInfo{Modifies:true},
    "hvals":            &CommandInfo{},

    //Lists
    //blpop
    //brpop
    "lindex":           &CommandInfo{},
    "linsert":          &CommandInfo{Modifies:true},
    "llen":             &CommandInfo{},
    "lpop":             &CommandInfo{Modifies:true},
    "lpush":            &CommandInfo{Modifies:true},
    "lpushx":           &CommandInfo{Modifies:true},
    "lrange":           &CommandInfo{},
    "lrem":             &CommandInfo{Modifies:true},
    "lset":             &CommandInfo{Modifies:true},
    "ltrim":            &CommandInfo{Modifies:true},
    "rpop":             &CommandInfo{Modifies:true},
    "rpush":            &CommandInfo{Modifies:true},
    "rpushx":           &CommandInfo{Modifies:true},

    //Sets
    "sadd":             &CommandInfo{Modifies:true},
    "scard":            &CommandInfo{},
    "sismember":        &CommandInfo{},
    "smembers":         &CommandInfo{},
    "spop":             &CommandInfo{Modifies:true},
    "srandmember":      &CommandInfo{},
    "srem":             &CommandInfo{Modifies:true},

    //Sorted Sets
    "zadd":             &CommandInfo{Modifies:true},
    "zcard":            &CommandInfo{},
    "zcount":           &CommandInfo{},
    "zincrby":          &CommandInfo{Modifies:true},
    "zrange":           &CommandInfo{},
    "zrangebyscore":    &CommandInfo{},
    "zrank":            &CommandInfo{},
    "zrem":             &CommandInfo{Modifies:true},
    "zremrangebyrank":  &CommandInfo{Modifies:true},
    "zremrangebyscore": &CommandInfo{Modifies:true},
    "zrevrange":        &CommandInfo{},
    "zrevrangebyscore": &CommandInfo{},
    "zrevrank":         &CommandInfo{},
    "zscore":           &CommandInfo{},

    //Monitors
    "madd":             &CommandInfo{IsCustom:true},
    "mrem":             &CommandInfo{IsCustom:true},

    //EKGs
    "eadd":             &CommandInfo{IsCustom:true,Modifies:true},
    "eaddq":            &CommandInfo{IsCustom:true,Modifies:true,IsQuiet:true},
    "erem":             &CommandInfo{IsCustom:true,Modifies:true},
    "eremq":            &CommandInfo{IsCustom:true,Modifies:true,IsQuiet:true},
    "emembers":         &CommandInfo{IsCustom:true},
}

// customCommandMap is a map of custom commands to their appropriate
// built-in funcions.
var customCommandMap = map[string]func(types.ConnId, *types.Payload)(interface{},error){

    //Monitors
    "madd":             custom.MAdd,
    "mrem":             custom.MRem,

    //EKGs
    "eadd":             custom.EAdd,
    "eaddq":            custom.EAdd,
    "erem":             custom.ERem,
    "eremq":            custom.ERem,
    "emembers":         custom.EMembers,

}

func GetCommandInfo(cmd *string) (*CommandInfo,bool) {
    cinfo,ok := commandMap[*cmd]
    return cinfo,ok
}
