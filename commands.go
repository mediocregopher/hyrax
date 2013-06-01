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
}

var commandMap = map[string]*CommandInfo{

    //Keys
    "del":              &CommandInfo{ WRITE_DIRECT, INT    },
    "exists":           &CommandInfo{ READ_DIRECT,  INT    },
    "expire":           &CommandInfo{ WRITE_DIRECT, INT    },
    "expireat":         &CommandInfo{ WRITE_DIRECT, INT    },
    "persist":          &CommandInfo{ WRITE_DIRECT, INT    },
    "pexpire":          &CommandInfo{ WRITE_DIRECT, INT    },
    "pexpireat":        &CommandInfo{ WRITE_DIRECT, INT    },
    "pttl":             &CommandInfo{ READ_DIRECT,  INT    },
    "ttl":              &CommandInfo{ READ_DIRECT,  INT    },
    "type":             &CommandInfo{ READ_DIRECT,  STRING },

    //Strings
    "append":           &CommandInfo{ WRITE_DIRECT, INT    },
    "bitcount":         &CommandInfo{ READ_DIRECT,  INT    },
    "decr":             &CommandInfo{ WRITE_DIRECT, INT    },
    "decrby":           &CommandInfo{ WRITE_DIRECT, INT    },
    "get":              &CommandInfo{ READ_DIRECT,  STRING },
    "getbit":           &CommandInfo{ READ_DIRECT,  INT    },
    "getrange":         &CommandInfo{ READ_DIRECT,  STRING },
    "getset":           &CommandInfo{ WRITE_DIRECT, STRING },
    "incr":             &CommandInfo{ WRITE_DIRECT, INT    },
    "incrby":           &CommandInfo{ WRITE_DIRECT, INT    },
    "incrbyfloat":      &CommandInfo{ WRITE_DIRECT, STRING },
    "mget":             &CommandInfo{ READ_DIRECT,  LIST   },
    "mset":             &CommandInfo{ WRITE_DIRECT, STRING },
    "msetnx":           &CommandInfo{ WRITE_DIRECT, INT    },
    "psetex":           &CommandInfo{ WRITE_DIRECT, STRING },
    "set":              &CommandInfo{ WRITE_DIRECT, STRING },
    "setbit":           &CommandInfo{ WRITE_DIRECT, INT    },
    "setex":            &CommandInfo{ WRITE_DIRECT, STRING },
    "setnx":            &CommandInfo{ WRITE_DIRECT, INT    },
    "setrange":         &CommandInfo{ WRITE_DIRECT, INT    },
    "strlen":           &CommandInfo{ READ_DIRECT,  INT    },

    //Hashes
    "hdel":             &CommandInfo{ WRITE_DIRECT, INT    },
    "hexists":          &CommandInfo{ READ_DIRECT,  INT    },
    "hget":             &CommandInfo{ READ_DIRECT,  STRING },
    "hgetall":          &CommandInfo{ READ_DIRECT,  LIST   },
    "hincrby":          &CommandInfo{ WRITE_DIRECT, INT    },
    "hincrbyfloat":     &CommandInfo{ WRITE_DIRECT, STRING },
    "hkeys":            &CommandInfo{ READ_DIRECT,  LIST   },
    "hlen":             &CommandInfo{ READ_DIRECT,  INT    },
    "hmget":            &CommandInfo{ READ_DIRECT,  LIST   },
    "hmset":            &CommandInfo{ WRITE_DIRECT, STRING },
    "hset":             &CommandInfo{ WRITE_DIRECT, INT    },
    "hsetnx":           &CommandInfo{ WRITE_DIRECT, INT    },
    "hvals":            &CommandInfo{ READ_DIRECT,  LIST   },

    //Lists
    //blpop
    //brpop
    "lindex":           &CommandInfo{ READ_DIRECT,  STRING },
    "linsert":          &CommandInfo{ WRITE_DIRECT, INT    },
    "llen":             &CommandInfo{ READ_DIRECT,  INT    },
    "lpop":             &CommandInfo{ WRITE_DIRECT, STRING },
    "lpush":            &CommandInfo{ WRITE_DIRECT, INT    },
    "lpushx":           &CommandInfo{ WRITE_DIRECT, INT    },
    "lrange":           &CommandInfo{ READ_DIRECT,  LIST   },
    "lrem":             &CommandInfo{ WRITE_DIRECT, INT    },
    "lset":             &CommandInfo{ WRITE_DIRECT, STRING },
    "ltrim":            &CommandInfo{ WRITE_DIRECT, STRING },
    "rpop":             &CommandInfo{ WRITE_DIRECT, STRING },
    "rpush":            &CommandInfo{ WRITE_DIRECT, INT    },
    "rpushx":           &CommandInfo{ WRITE_DIRECT, INT    },

    //Sets
    "sadd":             &CommandInfo{ WRITE_DIRECT, INT    },
    "scard":            &CommandInfo{ READ_DIRECT,  INT    },
    "sismember":        &CommandInfo{ READ_DIRECT,  INT    },
    "smembers":         &CommandInfo{ READ_DIRECT,  LIST   },
    "spop":             &CommandInfo{ WRITE_DIRECT, STRING },
    "srandmember":      &CommandInfo{ READ_DIRECT,  STRING },
    "srem":             &CommandInfo{ WRITE_DIRECT, INT    },

    //Sorted Sets
    "zadd":             &CommandInfo{ WRITE_DIRECT, INT    },
    "zcard":            &CommandInfo{ READ_DIRECT,  INT    },
    "zcount":           &CommandInfo{ READ_DIRECT,  INT    },
    "zincrby":          &CommandInfo{ WRITE_DIRECT, STRING },
    "zrange":           &CommandInfo{ READ_DIRECT,  LIST   },
    "zrangebyscore":    &CommandInfo{ READ_DIRECT,  LIST   },
    "zrank":            &CommandInfo{ READ_DIRECT,  INT    }, //TODO two different return values?
    "zrem":             &CommandInfo{ WRITE_DIRECT, INT    },
    "zremrangebyrank":  &CommandInfo{ WRITE_DIRECT, INT    },
    "zremrangebyscore": &CommandInfo{ WRITE_DIRECT, INT    },
    "zrevrange":        &CommandInfo{ READ_DIRECT,  LIST   },
    "zrevrangebyscore": &CommandInfo{ READ_DIRECT,  LIST   },
    "zrevrank":         &CommandInfo{ READ_DIRECT,  INT    }, //TODO two different return values?
    "zscore":           &CommandInfo{ READ_DIRECT,  STRING },
}
