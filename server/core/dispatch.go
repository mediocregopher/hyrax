package core

import (
	"errors"
	"github.com/mediocregopher/hyrax/server/auth"
	"github.com/mediocregopher/hyrax/server/dist"
	storage "github.com/mediocregopher/hyrax/server/storage-router"
	stypes "github.com/mediocregopher/hyrax/server/types"
	"github.com/mediocregopher/hyrax/types"
)

// RunCommand takes in a client id and a client command, figures out what type
// of command it is (builtin or direct) and routes the arguments to the
// appropriate function.
func RunCommand(
	cid stypes.ClientId,
	cmd *types.ClientCommand) *types.ClientReturn {

	var r interface{}
	var err error
	if storage.CommandFactory.DirectCommandAllowed(cmd.Command) {
		r, err = runDirectCommand(cid, cmd)
	} else if IsBuiltInCommand(cmd.Command) {
		r, err = runBuiltInCommand(cid, cmd)
	} else {
		err = errors.New("command not supported")
	}

	if err != nil {
		return &types.ClientReturn{Error: []byte(err.Error())}
	}

	return &types.ClientReturn{Return: r}
}

var autherr = errors.New("auth failed")

// runBuiltinCommand takes in a clientid and a builtin command struct and runs
// it, assuming auth checks out.
func runBuiltInCommand(
	cid stypes.ClientId,
	cmd *types.ClientCommand) (interface{}, error) {

	mods := CommandModifies(cmd.Command)
	adm := IsAdminCommand(cmd.Command)
	if mods || adm {
		ok, err := auth.Auth(cmd)
		if !ok {
			return nil, autherr
		} else if err != nil {
			return nil, err
		}
	}

	r, err := GetFunc(cmd.Command)(cid, cmd)

	if mods && !adm {
		dist.SendClientCommand(cmd)
	}

	return r, err
}

var directns = []byte("dir")

// runDirectCommand takes a client id and a client command, does authentication
// on the command if necessary, and runs the command directly on the correct
// storage unit for the command's key
func runDirectCommand(
	cid stypes.ClientId,
	cmd *types.ClientCommand) (interface{}, error) {

	mods := storage.CommandFactory.DirectCommandModifies(cmd.Command)
	if mods {
		ok, err := auth.Auth(cmd)
		if !ok {
			return nil, autherr
		} else if err != nil {
			return nil, err
		}
	}

	dcmd := storage.CommandFactory.DirectCommand(
		cmd.Command,
		storage.KeyMaker.Namespace(directns, cmd.StorageKey),
		cmd.Args,
	)

	r, err := storage.RoutedCmd(cmd.StorageKey, dcmd)

	if mods {
		dist.SendClientCommand(cmd)
	}

	return r, err
}

