package redis

import (
	"github.com/mediocregopher/hyrax/types"
	"github.com/mediocregopher/hyrax/server/storage-router/storage/command"
)

//These are for use by this and other modules so we don't have to re-allocate
//them everytime they get used
var ZADD = []byte("ZADD")
var ZREM = []byte("ZREM")
var ZMEMBERS = []byte("ZMEMBERS")
var ZCARD = []byte("ZCARD")
var ZRANGE = []byte("ZRANGE")
var ZRANGEBYSCORE = []byte("ZRANGEBYSCORE")
var ZSCORE = []byte("ZSCORE")

var SADD = []byte("SADD")
var SREM = []byte("SREM")
var SISMEMBER = []byte("SISMEMBER")
var SMEMBERS = []byte("SMEMBERS")
var SCARD = []byte("SCARD")

var MULTI = []byte("MULTI")
var EXEC = []byte("EXEC")


type RedisCommandFactory struct{}

func (r *RedisCommandFactory) createCmd(	
	cmd []byte,
	args ...interface{}) command.Command{

	return NewRedisCommand(cmd, args)
}

func (r *RedisCommandFactory) Transaction(
	cmds ...command.Command) command.Command {

	return NewRedisTransaction(cmds...)
}

func (r *RedisCommandFactory) DirectCommand(
	cmd []byte,
	key types.Byter,
	args []interface{}) command.Command {

	argv := append([]interface{}{key.Bytes()}, args...)
	return r.createCmd(cmd, argv...)
}

func (r *RedisCommandFactory) IdValueSetAdd(
	key types.Byter,
	id types.Uint64er,
	value types.Byter) command.Command {

	return r.createCmd(ZADD, key.Bytes(), id.Uint64(), value.Bytes())
}

// Theoretically this should remove a value that is (id,value) from the set, but
// redis sorted sets only support removing by (_,value), ignoring the id. The
// observable effect of this is that it will be possible for two connections to
// have the same value, one will just overwrite the other.
func (r *RedisCommandFactory) IdValueSetRem(
	key types.Byter,
	_ types.Uint64er,
	value types.Byter) command.Command {

	return r.createCmd(ZREM, key.Bytes(), value.Bytes())
}

// This will return an empty list if false, non-empty if true
func (r *RedisCommandFactory) IdValueSetIsIdMember(
	key types.Byter,
	id types.Uint64er) command.Command {

	return r.createCmd(ZRANGEBYSCORE, key.Bytes(), id.Uint64(), id.Uint64())
}

// This will return a nil if false, non-nil if true
func (r *RedisCommandFactory) IdValueSetIsValueMember(
	key types.Byter,
	value types.Byter) command.Command {

	return r.createCmd(ZSCORE, key.Bytes(), value.Bytes())
}

func (r *RedisCommandFactory) IdValueSetCard(key types.Byter) command.Command {
	return r.createCmd(ZCARD, key.Bytes())
}

func (r *RedisCommandFactory) IdValueSetMemberValues(
	key types.Byter) command.Command {

	return r.createCmd(ZRANGE, key.Bytes(), 0, -1)
}

func (r *RedisCommandFactory) GenericSetAdd(
	key types.Byter,
	value types.Byter) command.Command {

	return r.createCmd(SADD, key.Bytes(), value.Bytes())
}

func (r *RedisCommandFactory) GenericSetRem(
	key types.Byter,
	value types.Byter) command.Command {

	return r.createCmd(SREM, key.Bytes(), value.Bytes())
}

func (r *RedisCommandFactory) GenericSetIsMember(
	key types.Byter,
	value types.Byter) command.Command {

	return r.createCmd(SISMEMBER, key.Bytes(), value.Bytes())
}

func (r *RedisCommandFactory) GenericSetCard(key types.Byter) command.Command {
	return r.createCmd(SCARD, key.Bytes())
}

func (r *RedisCommandFactory) GenericSetMembers(
	key types.Byter) command.Command {
	return r.createCmd(SMEMBERS, key.Bytes())
}
