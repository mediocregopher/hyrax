package builtin

import (
	ctypes "github.com/mediocregopher/hyrax/types/client"
	stypes "github.com/mediocregopher/hyrax/server/types"
	"strings"
)

type BuiltInFunc func(
	stypes.ClientId,
	*ctypes.ClientCommand) (interface{}, error)

type builtInCommandInfo struct {
	Func BuiltInFunc
	Modifies bool
}

var builtInMap = map[string]*builtInCommandInfo{
	"madd":      &builtInCommandInfo{Func: MAdd, Modifies: true},
	"mrem":      &builtInCommandInfo{Func: MRem, Modifies: true},
	"eadd":      &builtInCommandInfo{Func: EAdd, Modifies: true},
	"erem":      &builtInCommandInfo{Func: ERem, Modifies: true},
	"emembers":  &builtInCommandInfo{Func: EMembers},
	"ecard":     &builtInCommandInfo{Func: ECard},
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

// GetFunc returns the function for a given builtin command, or nil if the
// command isn't a valid builtin command
func GetFunc(cmd []byte) BuiltInFunc {
	if cinfo, ok := getBuiltInCommandInfo(cmd); ok {
		return cinfo.Func
	}
	return nil
}
