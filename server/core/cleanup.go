package core

import (
	"github.com/mediocregopher/hyrax/server/dist"
	stypes "github.com/mediocregopher/hyrax/server/types"
	crouter "github.com/mediocregopher/hyrax/server/client-router"
	"github.com/mediocregopher/hyrax/types"
)

var closedCmd = []byte("eclose")

// ClientClosed takes care of all cleanup that's necessary when a client has
// closed
func ClientClosed(c stypes.Client) error {
	if err := CleanMons(c); err != nil {
		return err
	}

	ekgs, ids, err := EkgsForClient(c)
	if err != nil {
		return err
	}

	for i := range ekgs {
		cmd := &types.ClientCommand{
			Command:    closedCmd,
			StorageKey: ekgs[i],
			Id:         ids[i],
		}
		dist.SendClientCommand(cmd)
	}

	if err := CleanClientEkgsShort(ekgs, c); err != nil {
		return err
	}

	crouter.RemByClient(c)
	return crouter.UnsubscribeAll(c)
}
