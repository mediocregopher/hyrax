package redis

import (
	"github.com/mediocregopher/hyrax/server/storage-router/storage/command"
)

//These are for use by this and other modules so we don't have to re-allocate
//them everytime they get used
var DEL = []byte("DEL")

var HSET = []byte("HSET")
var HDEL = []byte("HDEL")
var HLEN = []byte("HLEN")
var HGETALL = []byte("HGETALL")
var HVALS = []byte("HVALS")

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
	cmd, key []byte,
	args []interface{}) command.Command {

	argv := append([]interface{}{key}, args...)
	return r.createCmd(cmd, argv...)
}

func (r *RedisCommandFactory) DirectCommandAllowed(cmd []byte) bool {
	_, ok := getCommandInfo(cmd)
	return ok
}

func (r *RedisCommandFactory) DirectCommandModifies(cmd []byte) bool {
	if cinfo, ok := getCommandInfo(cmd); ok {
		return cinfo.Modifies
	}

	return false
}

func (r *RedisCommandFactory) KeyValueSetAdd(
	key, innerkey, value []byte) command.Command {

	return r.createCmd(HSET, key, innerkey, value)
}

func (r *RedisCommandFactory) KeyValueSetRemByInnerKey(
	key, innerkey []byte) command.Command {

	return r.createCmd(HDEL, key,  innerkey)
}

func (r *RedisCommandFactory) KeyValueSetCard(key []byte) command.Command {
	return r.createCmd(HLEN, key)
}

func (r *RedisCommandFactory) KeyValueSetMembers(key []byte) command.Command {
	return r.createCmd(HGETALL, key)
}

func (r *RedisCommandFactory) KeyValueSetMemberValues(
	key []byte) command.Command {

	return r.createCmd(HVALS, key)
}

func (r *RedisCommandFactory) KeyValueSetDel(key []byte) command.Command {
	return r.createCmd(DEL, key)
}

func (r *RedisCommandFactory) GenericSetAdd(key, value []byte) command.Command {
	return r.createCmd(SADD, key, value)
}

func (r *RedisCommandFactory) GenericSetRem(key, value []byte) command.Command {
	return r.createCmd(SREM, key, value)
}

func (r *RedisCommandFactory) GenericSetIsMember(
	key, value []byte) command.Command {

	return r.createCmd(SISMEMBER, key, value)
}

func (r *RedisCommandFactory) GenericSetCard(key []byte) command.Command {
	return r.createCmd(SCARD, key)
}

func (r *RedisCommandFactory) GenericSetMembers(key []byte) command.Command {
	return r.createCmd(SMEMBERS, key)
}

func (r *RedisCommandFactory) GenericSetDel(key []byte) command.Command {
	return r.createCmd(DEL, key)
}
