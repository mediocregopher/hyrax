package core

import (
	"github.com/mediocregopher/hyrax/server/core/builtin"
	"github.com/mediocregopher/hyrax/server/core/keychanges"
	stypes "github.com/mediocregopher/hyrax/server/types"
	"github.com/mediocregopher/hyrax/types"
)

var closedCmd = "eclose"

// ClientClosed takes care of all cleanup that's necessary when a client has
// closed
func ClientClosed(c stypes.Client) error {
	ekgs, ids, err := builtin.EkgsForClient(c)
	if err != nil {
		return err
	}

	for i := range ekgs {
		cmd := &types.Action{
			Command:    closedCmd,
			StorageKey: ekgs[i],
			Id:         ids[i],
		}
		if err := keychanges.PubLocal(cmd); err != nil {
			return err
		}
	}

	if err := builtin.CleanClientEkgsShort(ekgs, c); err != nil {
		return err
	}

	return keychanges.UnsubscribeAll(c)
}
