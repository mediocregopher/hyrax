package client

import (
	"errors"
	"github.com/mediocregopher/hyrax/server/custom"
	"github.com/mediocregopher/hyrax/server/storage-router"
	"github.com/mediocregopher/hyrax/server/storage-router/storage"
	"github.com/mediocregopher/hyrax/types"
	ctypes "github.com/mediocregopher/hyrax/types/client"
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

// RunCommand takes in a client id and a client command, figures out what type
// of command it is (custom or direct) and routes the arguments to the
// appropriate function.
func RunCommand(
	cid stypes.ClientId,
	cmd *ctypes.ClientCommand) (interface{}, error) {

	if storage.CommandFactory.DirectCommandAllowed(cmd.Command) {
		return runDirectCommand(cid, cmd)
	} else if custom.IsCustomCommand(cmd.Command) {
		return runCustomCommand(cid, cmd)
	} else {
		return nil, errors.New("command not supported")
	}
}

// runCustomCommand takes in a clientid and a custom command struct and runs it,
// assuming auth checks out.
func runCustomCommand(
	cid stypes.ClientId,
	cmd *ctypes.ClientCommand) (interface{}, error) {

	if custom.CommandModifies(cmd.Command) {
		// Auth check
	}

	return custom.GetFunc(cmd.Command)(cid, cmd)
}

var directns = types.NewByter([]byte("dir"))

// runDirectCommand takes a client id and a client command, does authentication
// on the command if necessary, and runs the command directly on the correct
// storage unit for the command's key
func runDirectCommand(
	cid stypes.ClientId,
	cmd *ctypes.ClientCommand) (interface{}, error) {

	if storage.CommandFactory.DirectCommandModifies(cmd.Command) {
		// Auth check
	}

	dcmd := storage.CommandFactory.DirectCommand(
		cmd.Command,
		storage.KeyMaker.Namespace(directns, cmd.StorageKey),
		cmd.Args,
	)

	return router.RoutedCmd(cmd.StorageKey, dcmd)
}

// ClientClosed takes care of all cleanup that's necessary when a client has
// closed
func ClientClosed(cid stypes.ClientId) error {
	if err := custom.CleanMons(cid); err != nil {
		return err
	}

	// TODO: Send out push messages for ekgs
	if err := custom.CleanClientEkgs(cid); err != nil {
		return err
	}

	return nil
}
