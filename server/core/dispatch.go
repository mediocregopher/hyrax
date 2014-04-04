package core

import (
	"errors"
	"github.com/grooveshark/golib/gslog"
	"time"

	"github.com/mediocregopher/hyrax/server/auth"
	"github.com/mediocregopher/hyrax/server/config"
	"github.com/mediocregopher/hyrax/server/core/builtin"
	"github.com/mediocregopher/hyrax/server/core/keychanges"
	"github.com/mediocregopher/hyrax/server/listen"
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
	gslog.Infof("Connecting to datastore at %s", addr)
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

func init() {
	for i := 0; i < 10; i++ {
		go func() {
			for {
				select {
				case aw := <-listen.ActionWrapCh:
					go handleActionWrap(aw)
				case cc := <-listen.ClientClosedCh:
					go handleClientClosed(cc)
				}
			}
		}()
	}
}

func handleActionWrap(aw *listen.ActionWrap) {
	ar := RunAction(aw.Client, aw.Action)
	select {
	case aw.ActionReturnCh <- ar:
	case <-time.After(5 * time.Second):
		gslog.Warn("Client didn't read back ActionReturn")
	}
}

func handleClientClosed(cc *listen.ClientClosedWrap) {
	for i := 0; i < 3; i++ {
		if err := ClientClosed(cc.Client); err != nil {
			gslog.Errorf("calling ClientCleanup: %s", err)
			continue
		}
		close(cc.Ch)
		return
	}
}

// RunAction takes in a client and a client action, figures out what type of
// action it is (builtin or direct) and routes the arguments to the appropriate
// function.
func RunAction(c stypes.Client, cmd *types.Action) *types.ActionReturn {
	r, err := dispatchCommand(c, cmd)
	if err != nil {
		return types.NewActionReturn(err)
	}
	return types.NewActionReturn(r)
}

func dispatchCommand(c stypes.Client, cmd *types.Action) (interface{}, error) {

	var modifies, isAdmin func(string) bool
	var dispatch func(stypes.Client, *types.Action) (interface{}, error)
	if builtin.CommandIsBuiltIn(cmd.Command) {
		modifies = builtin.BuiltInCommandModifies
		isAdmin = builtin.BuiltInIsAdmin
		dispatch = builtin.GetBuiltInFunc(cmd.Command)
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
	if err == nil && mods && !adm {
		keychanges.PubLocal(cmd)
	}

	return r, err
}

// dispatchStorageCmd takes a client and a client command, and runs the command
// directly on the storage unit
func dispatchStorageCmd(
	c stypes.Client,
	cmd *types.Action) (interface{}, error) {

	args := make([]interface{}, 1, len(cmd.Args)+1)
	args[0] = cmd.StorageKey
	args = append(args, cmd.Args...)
	dcmd := storageUnit.NewCommand(cmd.Command, args...)
	return storageUnit.Cmd(dcmd)
}
