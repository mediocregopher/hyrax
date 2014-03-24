package core

import (
	"sync"

	stypes "github.com/mediocregopher/hyrax/server/types"
	"github.com/mediocregopher/hyrax/types"
)

// TODO two-way-map

// A mapping of keys being monitored to ClientIds
var monKeyToClientIds = map[string]map[uint64]bool{}

// A mapping of ClientIds to keys being monitored
var monClientIdToKeys = map[uint64]map[string]bool{}

// Lock which coordinates access to the mappings
var monLock sync.RWMutex

//MAdd adds the client's id to the set of clients that are monitoring the key
//(so it can receive alerts) and adds the key to the set of keys that the client
//is monitoring (so it can clean up)
func MAdd(cid stypes.ClientId, cmd *types.ClientCommand) (interface{}, error) {
	key := string(cmd.StorageKey)
	cidi := cid.Uint64()
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

// MRem removes the client's id from the set of clients that are monitoring the
// key, and removes the key from the set of keys that the client is monitoring
func MRem(cid stypes.ClientId, cmd *types.ClientCommand) (interface{}, error) {
	key := string(cmd.StorageKey)
	cidi := cid.Uint64()
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

// CleanMons takes in a client id and cleans up all of its monitors, and the set
// which keeps track of those monitors
func CleanMons(cid stypes.ClientId) error {
	cidi := cid.Uint64()
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

// ClientsForMon takes in a key and returns all the client ids on this node that
// are mon'ing that key
func ClientsForMon(keyb []byte) ([]stypes.ClientId, error) {
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
