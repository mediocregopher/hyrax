package client

import (
	"errors"
	"github.com/mediocregopher/hyrax/server/auth"
	"github.com/mediocregopher/hyrax/server/builtin"
	"github.com/mediocregopher/hyrax/server/dist"
	storage "github.com/mediocregopher/hyrax/server/storage-router"
	"github.com/mediocregopher/hyrax/types"
	stypes "github.com/mediocregopher/hyrax/server/types"
)

func init() {
	go idMakerSpin()
}

var idCh = make(chan stypes.ClientId)
func idMakerSpin() {
	for i := uint64(0) ;; i++ {
		// TODO do something with the error here (even though it'll never
		// happen)
		cid, _ := stypes.ClientIdFromUint64(i)
		idCh <- cid
	}
}

// NewClient returns a unique client id that a client can use to identify
// itself in later commands
func NewClient() stypes.ClientId {
	return <- idCh
}

// Client is an interface which must be implemented by clients to hyrax (go
// figure)
type Client interface {

	// ClientId returns the ClientId of a given client (again, go figure)
	ClientId() stypes.ClientId

	// CommandPushChannel returns a channel where commands that are to be pushed
	// to the client should be pushed on to
	CommandPushChannel() chan<- *types.ClientCommand
}

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
	} else if builtin.IsBuiltInCommand(cmd.Command) {
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

	mods := builtin.CommandModifies(cmd.Command)
	adm := builtin.IsAdminCommand(cmd.Command)
	if mods || adm {
		ok, err := auth.Auth(cmd)
		if !ok {
			return nil, autherr
		} else if err != nil {
			return nil, err
		}
	}

	r, err := builtin.GetFunc(cmd.Command)(cid, cmd)

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

var closedCmd = []byte("eclose")

// ClientClosed takes care of all cleanup that's necessary when a client has
// closed
func ClientClosed(cid stypes.ClientId) error {
	if err := builtin.CleanMons(cid); err != nil {
		return err
	}

	ekgs, ids, err := builtin.EkgsForClient(cid)
	if err != nil {
		return err
	}

	for i := range ekgs {
		cmd := &types.ClientCommand{
			Command: closedCmd,
			StorageKey: ekgs[i],
			Id: ids[i],
		}
		dist.SendClientCommand(cmd)
	}

	if err := builtin.CleanClientEkgsShort(ekgs, cid); err != nil {
		return err
	}

	return nil
}
