package custom

import (
	ctypes "github.com/mediocregopher/hyrax/types/client"
	stypes "github.com/mediocregopher/hyrax/server/types"
	"strings"
)

type CustomFunc func(stypes.ClientId, *ctypes.ClientCommand) (interface{}, error)

type customCommandInfo struct {
	Func CustomFunc
	Modifies bool
}

var customMap = map[string]*customCommandInfo{
	"madd":      &customCommandInfo{Func: MAdd, Modifies: true},
	"mrem":      &customCommandInfo{Func: MRem, Modifies: true},
	"eadd":      &customCommandInfo{Func: EAdd, Modifies: true},
	"erem":      &customCommandInfo{Func: ERem, Modifies: true},
	"emembers":  &customCommandInfo{Func: EMembers},
	"ecard":     &customCommandInfo{Func: ECard},
}

func getCustomCommandInfo(cmd []byte) (*customCommandInfo, bool) {
	cinfo, ok := customMap[strings.ToLower(string(cmd))]
	return cinfo, ok
}

// IsCustomCommand returns whether or not the given command is a valid custom
// one
func IsCustomCommand(cmd []byte) bool {
	_, ok := getCustomCommandInfo(cmd)
	return ok
}

// CommandModifies returns whether or not a given custom command modifies the
// backend state, or false if it's not a valid custom command
func CommandModifies(cmd []byte) bool {
	if cinfo, ok := getCustomCommandInfo(cmd); ok {
		return cinfo.Modifies
	}
	return false
}

// GetFunc returns the function for a given custom command, or nil if the
// command isn't a valid custom command
func GetFunc(cmd []byte) CustomFunc {
	if cinfo, ok := getCustomCommandInfo(cmd); ok {
		return cinfo.Func
	}
	return nil
}
