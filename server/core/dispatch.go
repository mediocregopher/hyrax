package core

import (
	"errors"
	"github.com/mediocregopher/hyrax/server/auth"
	"github.com/mediocregopher/hyrax/server/core/keychanges"
	storage "github.com/mediocregopher/hyrax/server/storage-router"
	stypes "github.com/mediocregopher/hyrax/server/types"
	"github.com/mediocregopher/hyrax/types"
)

// RunCommand takes in a client and a client command, figures out what type of
// command it is (builtin or direct) and routes the arguments to the appropriate
// function.
func RunCommand(c stypes.Client, cmd *types.ClientCommand) *types.ClientReturn {
	r, err := dispatchCommand(c, cmd)
	if err != nil {
		return &types.ClientReturn{Error: []byte(err.Error())}
	}

	return &types.ClientReturn{Return: r}
}

func dispatchCommand(c stypes.Client, cmd *types.ClientCommand) (interface{}, error) {
	mods := CommandModifies(cmd.Command)
	adm := IsAdminCommand(cmd.Command)
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
	cmd.Secret = nil

	var r interface{}
	var err error
	if storage.CommandFactory.DirectCommandAllowed(cmd.Command) {
		r, err = dispatchDirectCommand(c, cmd)
	} else if IsBuiltInCommand(cmd.Command) {
		r, err = GetFunc(cmd.Command)(c, cmd)
	} else {
		err = errors.New("command not supported")
	}

	if mods && !adm {
		keychanges.PubLocal(cmd)
	}

	return r, err
}

var directns = []byte("dir")

// runDirectCommand takes a client and a client command, does authentication on
// the command if necessary, and runs the command directly on the correct
// storage unit for the command's key
func dispatchDirectCommand(
	c stypes.Client,
	cmd *types.ClientCommand) (interface{}, error) {

	dcmd := storage.CommandFactory.DirectCommand(
		cmd.Command,
		storage.KeyMaker.Namespace(directns, cmd.StorageKey),
		cmd.Args,
	)

	return storage.RoutedCmd(cmd.StorageKey, dcmd)
}
