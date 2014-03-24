package core

import (
	"sync"

	stypes "github.com/mediocregopher/hyrax/server/types"
	"github.com/mediocregopher/hyrax/types"
	"github.com/mediocregopher/hyrax/server/core/keychanges"
)

// TODO two-way-map

// A mapping of keys being monitored to ClientIds
var monKeyToClientIds = map[string]map[uint64]bool{}

// A mapping of ClientIds to keys being monitored
var monClientIdToKeys = map[uint64]map[string]bool{}

// Lock which coordinates access to the mappings
var monLock sync.RWMutex

// MAll adds the client to the set of clients that are monitoring ALL keys. This
// hooks into a separate funtionality than the normal mon commands, so it will
// stack with them (aka, duplicate pushes if you also monitor individual keys)
func MAll(c stypes.Client, cmd *types.ClientCommand) (interface{}, error) {
	err := keychanges.AddClient(c)
	if err != nil {
		return nil, err
	}
	return []byte("OK"), nil
}

//MAdd adds the client to the set of clients that are monitoring the key (so it
//can receive alerts) and adds the key to the set of keys that the client is
//monitoring (so it can clean up)
func MAdd(c stypes.Client, cmd *types.ClientCommand) (interface{}, error) {
	key := string(cmd.StorageKey)
	cidi := c.ClientId().Uint64()
	monLock.Lock()
	defer monLock.Unlock()
	if clientIdsM, ok := monKeyToClientIds[key]; ok {
		clientIdsM[cidi] = true
	} else {
		monKeyToClientIds[key] = map[uint64]bool{cidi: true}
	}
	if keysM, ok := monClientIdToKeys[cidi]; ok {
		keysM[key] = true
	} else {
		monClientIdToKeys[cidi] = map[string]bool{key: true}
	}
	return []byte("OK"), nil
}

// MRem removes the client from the set of clients that are monitoring the key,
// and removes the key from the set of keys that the client is monitoring
func MRem(c stypes.Client, cmd *types.ClientCommand) (interface{}, error) {
	key := string(cmd.StorageKey)
	cidi := c.ClientId().Uint64()
	monLock.Lock()
	defer monLock.Unlock()

	// If the key isn't monitored, don't bother
	clientIdsM, ok := monKeyToClientIds[key]
	if !ok {
		return []byte("OK"), nil
	}

	if len(clientIdsM) == 1 {
		delete(monKeyToClientIds, key)
	}  else {
		delete(clientIdsM, cidi)
	}

	keysM := monClientIdToKeys[cidi]
	if len(keysM) == 1 {
		delete(monClientIdToKeys, cidi)
	} else {
		delete(keysM, key)
	}

	return []byte("OK"), nil
}

// CleanMons takes in a client and cleans up all of its monitors, and the set
// which keeps track of those monitors
func CleanMons(c stypes.Client) error {
	cidi := c.ClientId().Uint64()
	monLock.Lock()
	defer monLock.Unlock()


	// If the client isn't monitoring any keys, bail
	keysM, ok := monClientIdToKeys[cidi]
	if !ok {
		return nil
	}

	var clientIdsM map[uint64]bool
	for key := range keysM {
		clientIdsM = monKeyToClientIds[key]
		if len(clientIdsM) == 1 {
			delete(monKeyToClientIds, key)
		} else {
			delete(clientIdsM, cidi)
		}
	}

	return nil
}

// ClientIdsForMon takes in a key and returns all the clients on this node that
// are mon'ing that key
func ClientIdsForMon(keyb []byte) ([]stypes.ClientId, error) {
	key := string(keyb)
	monLock.RLock()
	defer monLock.RUnlock()

	clientIdsM, ok := monKeyToClientIds[key]
	if !ok {
		return []stypes.ClientId{}, nil
	}

	cids := make([]stypes.ClientId, 0, len(clientIdsM))
	for cidi := range clientIdsM {
		cid, err := stypes.ClientIdFromUint64(cidi)
		if err != nil {
			return nil, err
		}
		cids = append(cids, cid)
	}

	return cids, nil

}
