package dispatch

type ActionType int
const (
    READ_DIRECT ActionType = iota
    //READ_DIRECT_BLOCKING
    WRITE_DIRECT
    READ_CUSTOM
    WRITE_CUSTOM
)

//ReturnType
type ReturnType int
const (
    STRING ReturnType = iota
    INT
    FLOAT
    LIST
    MAP
)

type commandInfo struct {
    ActionType ActionType
    ReturnType ReturnType
}

var commandMap = map[string]*commandInfo{

    //Keys
    "del":              &commandInfo{ WRITE_DIRECT, INT    },
    "exists":           &commandInfo{ READ_DIRECT,  INT    },
    "expire":           &commandInfo{ WRITE_DIRECT, INT    },
    "expireat":         &commandInfo{ WRITE_DIRECT, INT    },
    "persist":          &commandInfo{ WRITE_DIRECT, INT    },
    "pexpire":          &commandInfo{ WRITE_DIRECT, INT    },
    "pexpireat":        &commandInfo{ WRITE_DIRECT, INT    },
    "pttl":             &commandInfo{ READ_DIRECT,  INT    },
    "ttl":              &commandInfo{ READ_DIRECT,  INT    },
    "type":             &commandInfo{ READ_DIRECT,  STRING },

    //Strings
    "append":           &commandInfo{ WRITE_DIRECT, INT    },
    "bitcount":         &commandInfo{ READ_DIRECT,  INT    },
    "decr":             &commandInfo{ WRITE_DIRECT, INT    },
    "decrby":           &commandInfo{ WRITE_DIRECT, INT    },
    "get":              &commandInfo{ READ_DIRECT,  STRING },
    "getbit":           &commandInfo{ READ_DIRECT,  INT    },
    "getrange":         &commandInfo{ READ_DIRECT,  STRING },
    "getset":           &commandInfo{ WRITE_DIRECT, STRING },
    "incr":             &commandInfo{ WRITE_DIRECT, INT    },
    "incrby":           &commandInfo{ WRITE_DIRECT, INT    },
    "incrbyfloat":      &commandInfo{ WRITE_DIRECT, STRING },
    "mget":             &commandInfo{ READ_DIRECT,  LIST   },
    "mset":             &commandInfo{ WRITE_DIRECT, STRING },
    "msetnx":           &commandInfo{ WRITE_DIRECT, INT    },
    "psetex":           &commandInfo{ WRITE_DIRECT, STRING },
    "set":              &commandInfo{ WRITE_DIRECT, STRING },
    "setbit":           &commandInfo{ WRITE_DIRECT, INT    },
    "setex":            &commandInfo{ WRITE_DIRECT, STRING },
    "setnx":            &commandInfo{ WRITE_DIRECT, INT    },
    "setrange":         &commandInfo{ WRITE_DIRECT, INT    },
    "strlen":           &commandInfo{ READ_DIRECT,  INT    },

    //Hashes
    "hdel":             &commandInfo{ WRITE_DIRECT, INT    },
    "hexists":          &commandInfo{ READ_DIRECT,  INT    },
    "hget":             &commandInfo{ READ_DIRECT,  STRING },
    "hgetall":          &commandInfo{ READ_DIRECT,  LIST   },
    "hincrby":          &commandInfo{ WRITE_DIRECT, INT    },
    "hincrbyfloat":     &commandInfo{ WRITE_DIRECT, STRING },
    "hkeys":            &commandInfo{ READ_DIRECT,  LIST   },
    "hlen":             &commandInfo{ READ_DIRECT,  INT    },
    "hmget":            &commandInfo{ READ_DIRECT,  LIST   },
    "hmset":            &commandInfo{ WRITE_DIRECT, STRING },
    "hset":             &commandInfo{ WRITE_DIRECT, INT    },
    "hsetnx":           &commandInfo{ WRITE_DIRECT, INT    },
    "hvals":            &commandInfo{ READ_DIRECT,  LIST   },

    //Lists
    //blpop
    //brpop
    "lindex":           &commandInfo{ READ_DIRECT,  STRING },
    "linsert":          &commandInfo{ WRITE_DIRECT, INT    },
    "llen":             &commandInfo{ READ_DIRECT,  INT    },
    "lpop":             &commandInfo{ WRITE_DIRECT, STRING },
    "lpush":            &commandInfo{ WRITE_DIRECT, INT    },
    "lpushx":           &commandInfo{ WRITE_DIRECT, INT    },
    "lrange":           &commandInfo{ READ_DIRECT,  LIST   },
    "lrem":             &commandInfo{ WRITE_DIRECT, INT    },
    "lset":             &commandInfo{ WRITE_DIRECT, STRING },
    "ltrim":            &commandInfo{ WRITE_DIRECT, STRING },
    "rpop":             &commandInfo{ WRITE_DIRECT, STRING },
    "rpush":            &commandInfo{ WRITE_DIRECT, INT    },
    "rpushx":           &commandInfo{ WRITE_DIRECT, INT    },

    //Sets
    "sadd":             &commandInfo{ WRITE_DIRECT, INT    },
    "scard":            &commandInfo{ READ_DIRECT,  INT    },
    "sismember":        &commandInfo{ READ_DIRECT,  INT    },
    "smembers":         &commandInfo{ READ_DIRECT,  LIST   },
    "spop":             &commandInfo{ WRITE_DIRECT, STRING },
    "srandmember":      &commandInfo{ READ_DIRECT,  STRING },
    "srem":             &commandInfo{ WRITE_DIRECT, INT    },

    //Sorted Sets
    "zadd":             &commandInfo{ WRITE_DIRECT, INT    },
    "zcard":            &commandInfo{ READ_DIRECT,  INT    },
    "zcount":           &commandInfo{ READ_DIRECT,  INT    },
    "zincrby":          &commandInfo{ WRITE_DIRECT, STRING },
    "zrange":           &commandInfo{ READ_DIRECT,  LIST   },
    "zrangebyscore":    &commandInfo{ READ_DIRECT,  LIST   },
    "zrank":            &commandInfo{ READ_DIRECT,  INT    }, //TODO two different return values?
    "zrem":             &commandInfo{ WRITE_DIRECT, INT    },
    "zremrangebyrank":  &commandInfo{ WRITE_DIRECT, INT    },
    "zremrangebyscore": &commandInfo{ WRITE_DIRECT, INT    },
    "zrevrange":        &commandInfo{ READ_DIRECT,  LIST   },
    "zrevrangebyscore": &commandInfo{ READ_DIRECT,  LIST   },
    "zrevrank":         &commandInfo{ READ_DIRECT,  INT    }, //TODO two different return values?
    "zscore":           &commandInfo{ READ_DIRECT,  STRING },
}
