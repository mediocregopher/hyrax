# Redis

Redis is the first supported storage backend for hyrax. Check the
[redis](http://redis.io) site for more information on the server itself. It's
the backend that's most conducive to hyrax, since hyrax was made with it in
mind.

Currently only a single redis node is usable, redis cluster will be supported in
the future (once it's actually released for production) and will be another
backend type.

## Configuration

When using Redis, the `storage-info` parameter for hyrax should take the form:

```
<host/ip>:<port>
```

With the default being:

```
localhost:6379
```

## Commands

The following commands are supported for being passed back to redis. Hyrax's
`key` field is passed back as the first argument to redis commands. For example,
a `SET FOO BAR` would look like:

```json
{"cmd":"SET","key":"FOO","args":["BAR"],"secret":"<hmac-sha1>"}
```

(Commands marked with `*` will require a `secret` field because they modify
state)

**Keys:**

* exists
* expire *
* expireat *
* persist *
* pexpire *
* pexpireat *
* pttl
* ttl
* type

**Strings:**

* append *
* bitcount
* decr *
* decrby *
* get
* getbit
* getrange
* getset *
* incr *
* incrby *
* incrbyfloat *
* psetex *
* set *
* setbit *
* setex *
* setnx *
* setrange *
* strlen

**Hashes:**

* hdel *
* hexists
* hget
* hgetall
* hincrby *
* hincrbyfloat *
* hkeys
* hlen
* hmget
* hset *
* hsetnx *
* hvals

**Lists:**

* lindex
* linsert *
* llen
* lpop *
* lpush *
* lpushx *
* lrange
* lrem *
* lset *
* ltrim *
* rpop *
* rpush *
* rpushx *

**Sets:**

* sadd *
* scard
* sismember
* smembers
* spop *
* srandmember
* srem *

**Sorted Sets:**

* zadd *
* zcard
* zcount
* zincrby *
* zrange
* zrangebyscore
* zrank
* zrem *
* zremrangebyrank *
* zremrangebyscore *
* zrevrange
* zrevrangebyscore
* zrevrank
* zscore
