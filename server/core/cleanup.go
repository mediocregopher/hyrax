package core

import (
	"github.com/mediocregopher/hyrax/server/dist"
	stypes "github.com/mediocregopher/hyrax/server/types"
	"github.com/mediocregopher/hyrax/types"
)

var closedCmd = []byte("eclose")

// ClientClosed takes care of all cleanup that's necessary when a client has
// closed
func ClientClosed(cid stypes.ClientId) error {
	if err := CleanMons(cid); err != nil {
		return err
	}

	ekgs, ids, err := EkgsForClient(cid)
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

	if err := CleanClientEkgsShort(ekgs, cid); err != nil {
		return err
	}

	return nil
}
