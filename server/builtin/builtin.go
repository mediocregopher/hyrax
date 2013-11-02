package builtin

import (
	"github.com/mediocregopher/hyrax/types"
	stypes "github.com/mediocregopher/hyrax/server/types"
	"strings"
)

type BuiltInFunc func(
	stypes.ClientId,
	*types.ClientCommand) (interface{}, error)

type builtInCommandInfo struct {
	Func BuiltInFunc
	Admin    bool
	Modifies bool
}

var builtInMap = map[string]*builtInCommandInfo{
	"madd":       &builtInCommandInfo{Func: MAdd, Modifies: true},
	"mrem":       &builtInCommandInfo{Func: MRem, Modifies: true},
	"eadd":       &builtInCommandInfo{Func: EAdd, Modifies: true},
	"erem":       &builtInCommandInfo{Func: ERem, Modifies: true},
	"emembers":   &builtInCommandInfo{Func: EMembers},
	"ecard":      &builtInCommandInfo{Func: ECard},

	"anodeadd":
		&builtInCommandInfo{
			Func: ANodeAdd,
			Modifies: true,
			Admin: true,
		},

	"anoderem":
		&builtInCommandInfo{
			Func: ANodeRem,
			Modifies: true,
			Admin: true,
		},

	"abucketset":
		&builtInCommandInfo{
			Func: ABucketSet,
			Modifies: true,
			Admin: true,
		},

	"abuckets":
		&builtInCommandInfo{
			Func: ABuckets,
			Admin: true,
		},

	"aglobalsecretadd":
		&builtInCommandInfo{
			Func: AGlobalSecretAdd,
			Modifies: true,
			Admin: true,
		},

	"aglobalsecretrem":
		&builtInCommandInfo{
			Func: AGlobalSecretRem,
			Modifies: true,
			Admin: true,
		},

	"aglobalsecrets":
		&builtInCommandInfo{
			Func: AGlobalSecrets,
			Admin: true,
		},

	"asecretadd":
		&builtInCommandInfo{
			Func: ASecretAdd,
			Modifies: true,
			Admin: true,
		},

	"asecretrem":
		&builtInCommandInfo{
			Func: ASecretRem,
			Modifies: true,
			Admin: true,
		},

	"asecrets":
		&builtInCommandInfo{
			Func: ASecrets,
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
