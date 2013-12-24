# Direct commands

Direct commands are commands which are forwarded straight through to the backend
data-store (redis). Hyrax allows you to use almost any redis command, with few
exceptions. You can find documentation for all redis commands at the [redis
website][redis].

The following are all redis commands that hyrax supports. Commands marked with a
"\*" are said to modify their key, and will generate a push message for clients
[monitoring][mon] that key.

Keys:
* exists
* expire *
* expireat *
* persist *
* pexpire *
* pexpireat *
* pttl
* ttl
* type

Strings:
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

Hashes:
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

Lists:
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

Sets:
* sadd *
* scard
* sismember
* smembers
* spop *
* srandmember
* srem *

Sorted Sets:
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


[redis]: http://redis.io/commands
[mon]: /doc/commands/mon.md
