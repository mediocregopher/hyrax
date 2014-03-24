package core

import (
	"github.com/mediocregopher/hyrax/server/config"
	storage "github.com/mediocregopher/hyrax/server/storage-router"
	stypes "github.com/mediocregopher/hyrax/server/types"
	"github.com/mediocregopher/hyrax/types"
	"strings"
)

// Some shortcuts which will be used by builtin functions
var keyMaker = storage.KeyMaker
var cmdFactory = storage.CommandFactory
var thisnode = &config.StorageAddr

type BuiltInFunc func(
	stypes.ClientId,
	*types.ClientCommand) (interface{}, error)

type builtInCommandInfo struct {
	Func     BuiltInFunc
	Admin    bool
	Modifies bool
}

var builtInMap = map[string]*builtInCommandInfo{
	"madd":     {Func: MAdd},
	"mrem":     {Func: MRem},
	"eadd":     {Func: EAdd, Modifies: true},
	"erem":     {Func: ERem, Modifies: true},
	"emembers": {Func: EMembers},
	"ecard":    {Func: ECard},

	"anodeadd": {
		Func:     ANodeAdd,
		Modifies: true,
		Admin:    true,
	},

	"anoderem": {
		Func:     ANodeRem,
		Modifies: true,
		Admin:    true,
	},

	"abucketset": {
		Func:     ABucketSet,
		Modifies: true,
		Admin:    true,
	},

	"abuckets": {
		Func:  ABuckets,
		Admin: true,
	},

	"aglobalsecretadd": {
		Func:     AGlobalSecretAdd,
		Modifies: true,
		Admin:    true,
	},

	"aglobalsecretrem": {
		Func:     AGlobalSecretRem,
		Modifies: true,
		Admin:    true,
	},

	"aglobalsecrets": {
		Func:  AGlobalSecrets,
		Admin: true,
	},

	"asecretadd": {
		Func:     ASecretAdd,
		Modifies: true,
		Admin:    true,
	},

	"asecretrem": {
		Func:     ASecretRem,
		Modifies: true,
		Admin:    true,
	},

	"asecrets": {
		Func:  ASecrets,
		Admin: true,
	},
}

func getBuiltInCommandInfo(cmd []byte) (*builtInCommandInfo, bool) {
	cinfo, ok := builtInMap[strings.ToLower(string(cmd))]
	return cinfo, ok
}

// IsBuiltInCommand returns whether or not the given command is a valid builtin
// one
func IsBuiltInCommand(cmd []byte) bool {
	_, ok := getBuiltInCommandInfo(cmd)
	return ok
}

// CommandModifies returns whether or not a given builtin command modifies the
// backend state, or false if it's not a valid builtin command
func CommandModifies(cmd []byte) bool {
	if cinfo, ok := getBuiltInCommandInfo(cmd); ok {
		return cinfo.Modifies
	}
	return false
}

// IsAdminCommand returns whether or not a given builtin command is an admin
// only command, or false if it's not a valid builtin command
func IsAdminCommand(cmd []byte) bool {
	if cinfo, ok := getBuiltInCommandInfo(cmd); ok {
		return cinfo.Admin
	}
	return false
}

// GetFunc returns the function for a given builtin command, or nil if the
// command isn't a valid builtin command
func GetFunc(cmd []byte) BuiltInFunc {
	if cinfo, ok := getBuiltInCommandInfo(cmd); ok {
		return cinfo.Func
	}
	return nil
}
