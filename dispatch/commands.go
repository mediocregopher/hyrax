package main

const (
    READ_DIRECT = iota
    //READ_DIRECT_BLOCKING
    WRITE_DIRECT
    READ_CUSTOM
    WRITE_CUSTOM
)

type CommandInfo struct {
    //Return type defined in json.go
    ActionType, ReturnType int
    MultipleKeys bool
}

var commandMap = map[string]*CommandInfo{

    //Keys
    "del":              &CommandInfo{ WRITE_DIRECT, INT,    true  },
    "exists":           &CommandInfo{ READ_DIRECT,  INT,    false },
    "expire":           &CommandInfo{ WRITE_DIRECT, INT,    false },
    "expireat":         &CommandInfo{ WRITE_DIRECT, INT,    false },
    "persist":          &CommandInfo{ WRITE_DIRECT, INT,    false },
    "pexpire":          &CommandInfo{ WRITE_DIRECT, INT,    false },
    "pexpireat":        &CommandInfo{ WRITE_DIRECT, INT,    false },
    "pttl":             &CommandInfo{ READ_DIRECT,  INT,    false },
    "ttl":              &CommandInfo{ READ_DIRECT,  INT,    false },
    "type":             &CommandInfo{ READ_DIRECT,  STRING, false },

    //Strings
    "append":           &CommandInfo{ WRITE_DIRECT, INT,    false },
    "bitcount":         &CommandInfo{ READ_DIRECT,  INT,    false },
    "decr":             &CommandInfo{ WRITE_DIRECT, INT,    false },
    "decrby":           &CommandInfo{ WRITE_DIRECT, INT,    false },
    "get":              &CommandInfo{ READ_DIRECT,  STRING, false },
    "getbit":           &CommandInfo{ READ_DIRECT,  INT,    false },
    "getrange":         &CommandInfo{ READ_DIRECT,  STRING, false },
    "getset":           &CommandInfo{ WRITE_DIRECT, STRING, false },
    "incr":             &CommandInfo{ WRITE_DIRECT, INT,    false },
    "incrby":           &CommandInfo{ WRITE_DIRECT, INT,    false },
    "incrbyfloat":      &CommandInfo{ WRITE_DIRECT, STRING, false },
    "mget":             &CommandInfo{ READ_DIRECT,  LIST,   true  },
    "mset":             &CommandInfo{ WRITE_DIRECT, STRING, true  },
    "msetnx":           &CommandInfo{ WRITE_DIRECT, INT,    true  },
    "psetex":           &CommandInfo{ WRITE_DIRECT, STRING, false },
    "set":              &CommandInfo{ WRITE_DIRECT, STRING, false },
    "setbit":           &CommandInfo{ WRITE_DIRECT, INT,    false },
    "setex":            &CommandInfo{ WRITE_DIRECT, STRING, false },
    "setnx":            &CommandInfo{ WRITE_DIRECT, INT,    false },
    "setrange":         &CommandInfo{ WRITE_DIRECT, INT,    false },
    "strlen":           &CommandInfo{ READ_DIRECT,  INT,    false },

    //Hashes
    "hdel":             &CommandInfo{ WRITE_DIRECT, INT,    false },
    "hexists":          &CommandInfo{ READ_DIRECT,  INT,    false },
    "hget":             &CommandInfo{ READ_DIRECT,  STRING, false },
    "hgetall":          &CommandInfo{ READ_DIRECT,  LIST,   false },
    "hincrby":          &CommandInfo{ WRITE_DIRECT, INT,    false },
    "hincrbyfloat":     &CommandInfo{ WRITE_DIRECT, STRING, false },
    "hkeys":            &CommandInfo{ READ_DIRECT,  LIST,   false },
    "hlen":             &CommandInfo{ READ_DIRECT,  INT,    false },
    "hmget":            &CommandInfo{ READ_DIRECT,  LIST,   false },
    "hmset":            &CommandInfo{ WRITE_DIRECT, STRING, true  },
    "hset":             &CommandInfo{ WRITE_DIRECT, INT,    false },
    "hsetnx":           &CommandInfo{ WRITE_DIRECT, INT,    false },
    "hvals":            &CommandInfo{ READ_DIRECT,  LIST,   false },

    //Lists
    //blpop
    //brpop
    "lindex":           &CommandInfo{ READ_DIRECT,  STRING, false },
    "linsert":          &CommandInfo{ WRITE_DIRECT, INT,    false },
    "llen":             &CommandInfo{ READ_DIRECT,  INT,    false },
    "lpop":             &CommandInfo{ WRITE_DIRECT, STRING, false },
    "lpush":            &CommandInfo{ WRITE_DIRECT, INT,    false },
    "lpushx":           &CommandInfo{ WRITE_DIRECT, INT,    false },
    "lrange":           &CommandInfo{ READ_DIRECT,  LIST,   false },
    "lrem":             &CommandInfo{ WRITE_DIRECT, INT,    false },
    "lset":             &CommandInfo{ WRITE_DIRECT, STRING, false },
    "ltrim":            &CommandInfo{ WRITE_DIRECT, STRING, false },
    "rpop":             &CommandInfo{ WRITE_DIRECT, STRING, false },
    "rpush":            &CommandInfo{ WRITE_DIRECT, INT,    false },
    "rpushx":           &CommandInfo{ WRITE_DIRECT, INT,    false },

    //Sets
    "sadd":             &CommandInfo{ WRITE_DIRECT, INT,    false },
    "scard":            &CommandInfo{ READ_DIRECT,  INT,    false },
    "sismember":        &CommandInfo{ READ_DIRECT,  INT,    false },
    "smembers":         &CommandInfo{ READ_DIRECT,  LIST,   false },
    "spop":             &CommandInfo{ WRITE_DIRECT, STRING, false },
    "srandmember":      &CommandInfo{ READ_DIRECT,  STRING, false },
    "srem":             &CommandInfo{ WRITE_DIRECT, INT,    false },

    //Sorted Sets
    "zadd":             &CommandInfo{ WRITE_DIRECT, INT,    false },
    "zcard":            &CommandInfo{ READ_DIRECT,  INT,    false },
    "zcount":           &CommandInfo{ READ_DIRECT,  INT,    false },
    "zincrby":          &CommandInfo{ WRITE_DIRECT, STRING, false },
    "zrange":           &CommandInfo{ READ_DIRECT,  LIST,   false },
    "zrangebyscore":    &CommandInfo{ READ_DIRECT,  LIST,   false },
    "zrank":            &CommandInfo{ READ_DIRECT,  INT,    false }, //TODO two different return values?
    "zrem":             &CommandInfo{ WRITE_DIRECT, INT,    false },
    "zremrangebyrank":  &CommandInfo{ WRITE_DIRECT, INT,    false },
    "zremrangebyscore": &CommandInfo{ WRITE_DIRECT, INT,    false },
    "zrevrange":        &CommandInfo{ READ_DIRECT,  LIST,   false },
    "zrevrangebyscore": &CommandInfo{ READ_DIRECT,  LIST,   false },
    "zrevrank":         &CommandInfo{ READ_DIRECT,  INT,    false }, //TODO two different return values?
    "zscore":           &CommandInfo{ READ_DIRECT,  STRING, false },
}
