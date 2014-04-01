package core

import (
	"errors"

	"github.com/mediocregopher/hyrax/server/auth"
	"github.com/mediocregopher/hyrax/server/config"
	"github.com/mediocregopher/hyrax/server/core/keychanges"
	"github.com/mediocregopher/hyrax/server/storage"
	"github.com/mediocregopher/hyrax/server/storage/redis"
	stypes "github.com/mediocregopher/hyrax/server/types"
	"github.com/mediocregopher/hyrax/types"
)

// The set of connections into the actual data store
var storageUnit *storage.StorageUnit

// The number of connections in the storage unit
const UNITSIZE = 10

func SetupStorage() error {
	// We assume redis for now since it's the only available type
	addr := config.StorageInfo
	sucs := make([]storage.Storage, UNITSIZE)
	for i := range sucs {
		sucs[i] = redis.New()
	}
	su, err := storage.NewStorageUnit(sucs, "tcp", addr)
	if err != nil {
		return err
	}
	storageUnit = su
	return nil
}

// RunCommand takes in a client and a client command, figures out what type of
// command it is (builtin or direct) and routes the arguments to the appropriate
// function.
func RunCommand(c stypes.Client, cmd *types.ClientCommand) *types.ClientReturn {
	r, err := dispatchCommand(c, cmd)
	if err != nil {
		return types.ErrorReturn(err)
	}

	return &types.ClientReturn{Return: r}
}

func dispatchCommand(c stypes.Client, cmd *types.ClientCommand) (interface{}, error) {

	var modifies, isAdmin func(string) bool
	var dispatch func(stypes.Client, *types.ClientCommand) (interface{}, error)
	if CommandIsBuiltIn(cmd.Command) {
		modifies = BuiltInCommandModifies
		isAdmin = BuiltInIsAdmin
		dispatch = GetBuiltInFunc(cmd.Command)
	} else if storageUnit.CommandAllowed(cmd.Command) {
		modifies = storageUnit.CommandModifies
		isAdmin = storageUnit.CommandIsAdmin
		dispatch = dispatchStorageCmd
	} else {
		return nil, errors.New("command not supported")
	}

	mods := modifies(cmd.Command)
	adm := isAdmin(cmd.Command)
	if mods || adm {
		ok, err := auth.Auth(cmd)
		if !ok {
			return nil, errors.New("auth failed")
		} else if err != nil {
			return nil, err
		}
	}
	// Before this cmd can get sent outside this go-routine we want to make sure
	// the secret is cleared
	cmd.Secret = ""

	r, err := dispatch(c, cmd)
	if err != nil && mods && !adm {
		keychanges.PubLocal(cmd)
	}

	return r, err
}

// dispatchStorageCmd takes a client and a client command, and runs the command
// directly on the storage unit
func dispatchStorageCmd(
	c stypes.Client,
	cmd *types.ClientCommand) (interface{}, error) {

	args := make([]interface{}, 1, len(cmd.Args)+1)
	args[0] = cmd.StorageKey
	args = append(args, cmd.Args...)
	dcmd := storageUnit.NewCommand(cmd.Command, args)
	return storageUnit.Cmd(dcmd)
}
