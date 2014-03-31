package core

import (
	stypes "github.com/mediocregopher/hyrax/server/types"
	"github.com/mediocregopher/hyrax/types"
	"strings"
)

var OK = []byte("OK")

type BuiltInFunc func(
	stypes.Client,
	*types.ClientCommand) (interface{}, error)

type builtInCommandInfo struct {
	Func     BuiltInFunc
	Admin    bool
	Modifies bool
}

var builtInMap = map[string]*builtInCommandInfo{
	"mglobal":  {Func: MGlobal, Admin: true},
	"mlocal":   {Func: MLocal, Admin: true},
	"madd":     {Func: MAdd},
	"mrem":     {Func: MRem},
	"eadd":     {Func: EAdd, Modifies: true},
	"erem":     {Func: ERem, Modifies: true},
	"emembers": {Func: EMembers},
	"ecard":    {Func: ECard},

	"alistentome": {Func: AListenToMe, Admin: true},
	"aignoreme":   {Func: AIgnoreMe, Admin: true},

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

	//"asecretadd": {
	//	Func:     ASecretAdd,
	//	Modifies: true,
	//	Admin:    true,
	//},

	//"asecretrem": {
	//	Func:     ASecretRem,
	//	Modifies: true,
	//	Admin:    true,
	//},

	//"asecrets": {
	//	Func:  ASecrets,
	//	Admin: true,
	//},
}

func getBuiltInCommandInfo(cmd []byte) (*builtInCommandInfo, bool) {
	cinfo, ok := builtInMap[strings.ToLower(string(cmd))]
	return cinfo, ok
}

// CommandIsBuiltIn returns whether or not the given command is a valid builtin
// one
func CommandIsBuiltIn(cmd []byte) bool {
	_, ok := getBuiltInCommandInfo(cmd)
	return ok
}

// BuiltInCommandModfies returns whether or not a given builtin command modifies
// the backend state, or false if it's not a valid builtin command
func BuiltInCommandModifies(cmd []byte) bool {
	if cinfo, ok := getBuiltInCommandInfo(cmd); ok {
		return cinfo.Modifies
	}
	return false
}

// BuiltInIsAdmin returns whether or not a given builtin command is an admin
// only command, or false if it's not a valid builtin command
func BuiltInIsAdmin(cmd []byte) bool {
	if cinfo, ok := getBuiltInCommandInfo(cmd); ok {
		return cinfo.Admin
	}
	return false
}

// GetBuiltInFunc returns the function for a given builtin command, or nil if
// the command isn't a valid builtin command
func GetBuiltInFunc(cmd []byte) BuiltInFunc {
	if cinfo, ok := getBuiltInCommandInfo(cmd); ok {
		return cinfo.Func
	}
	return nil
}
