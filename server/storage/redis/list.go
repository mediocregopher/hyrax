package redis

import (
	"strings"
)

// CommandInfo is a struct which is tied to a command, and describes various
// properties of the command. All properties are false by default.
type CommandInfo struct {
	Modifies bool
}

// commandMap is a map of commands to their info structs
var commandMap = map[string]*CommandInfo{

	//Keys
	"exists":    {},
	"expire":    {Modifies: true},
	"expireat":  {Modifies: true},
	"persist":   {Modifies: true},
	"pexpire":   {Modifies: true},
	"pexpireat": {Modifies: true},
	"pttl":      {},
	"ttl":       {},
	"type":      {},

	//Strings
	"append":      {Modifies: true},
	"bitcount":    {},
	"decr":        {Modifies: true},
	"decrby":      {Modifies: true},
	"get":         {},
	"getbit":      {},
	"getrange":    {},
	"getset":      {Modifies: true},
	"incr":        {Modifies: true},
	"incrby":      {Modifies: true},
	"incrbyfloat": {Modifies: true},
	"psetex":      {Modifies: true},
	"set":         {Modifies: true},
	"setbit":      {Modifies: true},
	"setex":       {Modifies: true},
	"setnx":       {Modifies: true},
	"setrange":    {Modifies: true},
	"strlen":      {},

	//Hashes
	"hdel":         {Modifies: true},
	"hexists":      {},
	"hget":         {},
	"hgetall":      {}, // Returns map
	"hincrby":      {Modifies: true},
	"hincrbyfloat": {Modifies: true},
	"hkeys":        {},
	"hlen":         {},
	"hmget":        {},
	"hset":         {Modifies: true},
	"hsetnx":       {Modifies: true},
	"hvals":        {},

	//Lists
	//blpop
	//brpop
	"lindex":  {},
	"linsert": {Modifies: true},
	"llen":    {},
	"lpop":    {Modifies: true},
	"lpush":   {Modifies: true},
	"lpushx":  {Modifies: true},
	"lrange":  {},
	"lrem":    {Modifies: true},
	"lset":    {Modifies: true},
	"ltrim":   {Modifies: true},
	"rpop":    {Modifies: true},
	"rpush":   {Modifies: true},
	"rpushx":  {Modifies: true},

	//Sets
	"sadd":        {Modifies: true},
	"scard":       {},
	"sismember":   {},
	"smembers":    {},
	"spop":        {Modifies: true},
	"srandmember": {},
	"srem":        {Modifies: true},

	//Sorted Sets
	"zadd":             {Modifies: true},
	"zcard":            {},
	"zcount":           {},
	"zincrby":          {Modifies: true},
	"zrange":           {},
	"zrangebyscore":    {},
	"zrank":            {},
	"zrem":             {Modifies: true},
	"zremrangebyrank":  {Modifies: true},
	"zremrangebyscore": {Modifies: true},
	"zrevrange":        {},
	"zrevrangebyscore": {},
	"zrevrank":         {},
	"zscore":           {},
}

func getCommandInfo(cmd []byte) (*CommandInfo, bool) {
	cinfo, ok := commandMap[strings.ToLower(string(cmd))]
	return cinfo, ok
}
