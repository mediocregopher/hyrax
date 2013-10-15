package redis

import (
	"github.com/mediocregopher/hyrax/types"
	"github.com/mediocregopher/hyrax/server/storage-router/storage/command"
)

//These are for use by this and other modules so we don't have to re-allocate
//them everytime they get used
var DEL = []byte("DEL")

var HSET = []byte("HSET")
var HDEL = []byte("HDEL")
var HLEN = []byte("HLEN")
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
	cmd []byte,
	key types.Byter,
	args []interface{}) command.Command {

	argv := append([]interface{}{key.Bytes()}, args...)
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
	key types.Byter,
	innerkey types.Byter,
	value types.Byter) command.Command {

	return r.createCmd(HSET, key.Bytes(), innerkey.Bytes(), value.Bytes())
}

func (r *RedisCommandFactory) KeyValueSetRemByInnerKey(
	key types.Byter,
	innerkey types.Byter) command.Command {

	return r.createCmd(HDEL, key.Bytes(),  innerkey.Bytes())
}

func (r *RedisCommandFactory) KeyValueSetCard(key types.Byter) command.Command {
	return r.createCmd(HLEN, key.Bytes())
}

func (r *RedisCommandFactory) KeyValueSetMemberValues(
	key types.Byter) command.Command {

	return r.createCmd(HVALS, key.Bytes())
}

func (r *RedisCommandFactory) KeyValueSetDel(key types.Byter) command.Command {
	return r.createCmd(DEL, key.Bytes())
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

func (r *RedisCommandFactory) GenericSetDel(key types.Byter) command.Command {
	return r.createCmd(DEL, key.Bytes())
}
