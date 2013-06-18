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
    MultipleKeys bool
}

var commandMap = map[string]*commandInfo{

    //Keys
    "exists":           &commandInfo{ READ_DIRECT,  INT,    false },
    "expire":           &commandInfo{ WRITE_DIRECT, INT,    false },
    "expireat":         &commandInfo{ WRITE_DIRECT, INT,    false },
    "persist":          &commandInfo{ WRITE_DIRECT, INT,    false },
    "pexpire":          &commandInfo{ WRITE_DIRECT, INT,    false },
    "pexpireat":        &commandInfo{ WRITE_DIRECT, INT,    false },
    "pttl":             &commandInfo{ READ_DIRECT,  INT,    false },
    "ttl":              &commandInfo{ READ_DIRECT,  INT,    false },
    "type":             &commandInfo{ READ_DIRECT,  STRING, false },

    //Strings
    "append":           &commandInfo{ WRITE_DIRECT, INT,    false },
    "bitcount":         &commandInfo{ READ_DIRECT,  INT,    false },
    "decr":             &commandInfo{ WRITE_DIRECT, INT,    false },
    "decrby":           &commandInfo{ WRITE_DIRECT, INT,    false },
    "get":              &commandInfo{ READ_DIRECT,  STRING, false },
    "getbit":           &commandInfo{ READ_DIRECT,  INT,    false },
    "getrange":         &commandInfo{ READ_DIRECT,  STRING, false },
    "getset":           &commandInfo{ WRITE_DIRECT, STRING, false },
    "incr":             &commandInfo{ WRITE_DIRECT, INT,    false },
    "incrby":           &commandInfo{ WRITE_DIRECT, INT,    false },
    "incrbyfloat":      &commandInfo{ WRITE_DIRECT, STRING, false },
    "psetex":           &commandInfo{ WRITE_DIRECT, STRING, false },
    "set":              &commandInfo{ WRITE_DIRECT, STRING, false },
    "setbit":           &commandInfo{ WRITE_DIRECT, INT,    false },
    "setex":            &commandInfo{ WRITE_DIRECT, STRING, false },
    "setnx":            &commandInfo{ WRITE_DIRECT, INT,    false },
    "setrange":         &commandInfo{ WRITE_DIRECT, INT,    false },
    "strlen":           &commandInfo{ READ_DIRECT,  INT,    false },

    //Hashes
    "hdel":             &commandInfo{ WRITE_DIRECT, INT,    false },
    "hexists":          &commandInfo{ READ_DIRECT,  INT,    false },
    "hget":             &commandInfo{ READ_DIRECT,  STRING, false },
    "hgetall":          &commandInfo{ READ_DIRECT,  MAP,    false },
    "hincrby":          &commandInfo{ WRITE_DIRECT, INT,    false },
    "hincrbyfloat":     &commandInfo{ WRITE_DIRECT, STRING, false },
    "hkeys":            &commandInfo{ READ_DIRECT,  LIST,   false },
    "hlen":             &commandInfo{ READ_DIRECT,  INT,    false },
    "hmget":            &commandInfo{ READ_DIRECT,  LIST,   false },
    "hset":             &commandInfo{ WRITE_DIRECT, INT,    false },
    "hsetnx":           &commandInfo{ WRITE_DIRECT, INT,    false },
    "hvals":            &commandInfo{ READ_DIRECT,  LIST,   false },

    //Lists
    //blpop
    //brpop
    "lindex":           &commandInfo{ READ_DIRECT,  STRING, false },
    "linsert":          &commandInfo{ WRITE_DIRECT, INT,    false },
    "llen":             &commandInfo{ READ_DIRECT,  INT,    false },
    "lpop":             &commandInfo{ WRITE_DIRECT, STRING, false },
    "lpush":            &commandInfo{ WRITE_DIRECT, INT,    false },
    "lpushx":           &commandInfo{ WRITE_DIRECT, INT,    false },
    "lrange":           &commandInfo{ READ_DIRECT,  LIST,   false },
    "lrem":             &commandInfo{ WRITE_DIRECT, INT,    false },
    "lset":             &commandInfo{ WRITE_DIRECT, STRING, false },
    "ltrim":            &commandInfo{ WRITE_DIRECT, STRING, false },
    "rpop":             &commandInfo{ WRITE_DIRECT, STRING, false },
    "rpush":            &commandInfo{ WRITE_DIRECT, INT,    false },
    "rpushx":           &commandInfo{ WRITE_DIRECT, INT,    false },

    //Sets
    "sadd":             &commandInfo{ WRITE_DIRECT, INT,    false },
    "scard":            &commandInfo{ READ_DIRECT,  INT,    false },
    "sismember":        &commandInfo{ READ_DIRECT,  INT,    false },
    "smembers":         &commandInfo{ READ_DIRECT,  LIST,   false },
    "spop":             &commandInfo{ WRITE_DIRECT, STRING, false },
    "srandmember":      &commandInfo{ READ_DIRECT,  STRING, false },
    "srem":             &commandInfo{ WRITE_DIRECT, INT,    false },

    //Sorted Sets
    "zadd":             &commandInfo{ WRITE_DIRECT, INT,    false },
    "zcard":            &commandInfo{ READ_DIRECT,  INT,    false },
    "zcount":           &commandInfo{ READ_DIRECT,  INT,    false },
    "zincrby":          &commandInfo{ WRITE_DIRECT, STRING, false },
    "zrange":           &commandInfo{ READ_DIRECT,  LIST,   false },
    "zrangebyscore":    &commandInfo{ READ_DIRECT,  LIST,   false },
    "zrank":            &commandInfo{ READ_DIRECT,  INT,    false }, //TODO two different return values?
    "zrem":             &commandInfo{ WRITE_DIRECT, INT,    false },
    "zremrangebyrank":  &commandInfo{ WRITE_DIRECT, INT,    false },
    "zremrangebyscore": &commandInfo{ WRITE_DIRECT, INT,    false },
    "zrevrange":        &commandInfo{ READ_DIRECT,  LIST,   false },
    "zrevrangebyscore": &commandInfo{ READ_DIRECT,  LIST,   false },
    "zrevrank":         &commandInfo{ READ_DIRECT,  INT,    false }, //TODO two different return values?
    "zscore":           &commandInfo{ READ_DIRECT,  STRING, false },
}
