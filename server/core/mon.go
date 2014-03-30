package core

import (
	"github.com/mediocregopher/hyrax/server/core/keychanges"
	stypes "github.com/mediocregopher/hyrax/server/types"
	"github.com/mediocregopher/hyrax/server/dist2"
	"github.com/mediocregopher/hyrax/types"
)

var GlobalManager = dist2.New("MGLOBAL")
var LocalManager = dist2.New("MLOCAL")

// MGlobal adds the client to the set of clients that are monitoring key changes
// happening all over the cluster. This hooks into a separate funtionality than
// the normal mon commands, so it will stack with them (aka, duplicate pushes if
// you also monitor individual keys)
func MGlobal(c stypes.Client, cmd *types.ClientCommand) (interface{}, error) {
	return OK, keychanges.SubscribeGlobal(c)
}

// MLocal adds the client to the set of clients that are monitoring key changes
// happening on this node of the cluster. This hooks into a separate
// funtionality than the normal mon commands, so it will stack with them (aka,
// duplicate pushes if you also monitor individual keys)
func MLocal(c stypes.Client, cmd *types.ClientCommand) (interface{}, error) {
	return OK, keychanges.SubscribeLocal(c)
}

//MAdd adds the client to the set of clients that are monitoring the key (so it
//can receive alerts) and adds the key to the set of keys that the client is
//monitoring (so it can clean up)
func MAdd(c stypes.Client, cmd *types.ClientCommand) (interface{}, error) {
	return OK, keychanges.Mon(c, cmd.StorageKey)
}

// MRem removes the client from the set of clients that are monitoring the key,
// and removes the key from the set of keys that the client is monitoring
func MRem(c stypes.Client, cmd *types.ClientCommand) (interface{}, error) {
	return OK, keychanges.Unmon(c, cmd.StorageKey)
}
