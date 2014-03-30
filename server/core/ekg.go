package core

import (
	"sync"

	stypes "github.com/mediocregopher/hyrax/server/types"
	"github.com/mediocregopher/hyrax/types"
)

// A mapping of ekgs to ClientIds and their names
var ekgKeyToClientIdsNames = map[string]map[uint64][]byte{}

// A mapping of ClientIds to the ekgs the client is on and the names it's using
// for each ekg
var ekgClientIdToKeysNames = map[uint64]map[string][]byte{}

// Lock which coordinates access to the mappings
var ekgLock sync.RWMutex

// EAdd adds the client to an ekg's set of things it's watching, and adds the
// ekg's information to the client's set of ekgs its hooked up to
func EAdd(c stypes.Client, cmd *types.ClientCommand) (interface{}, error) {
	key := string(cmd.StorageKey)
	cidi := c.ClientId().Uint64()
	name := cmd.Id
	ekgLock.Lock()
	defer ekgLock.Unlock()

	if clientIdsM, ok := ekgKeyToClientIdsNames[key]; ok {
		clientIdsM[cidi] = name
	} else {
		ekgKeyToClientIdsNames[key] = map[uint64][]byte{cidi: name}
	}
	if keysM, ok := ekgClientIdToKeysNames[cidi]; ok {
		keysM[key] = name
	} else {
		ekgClientIdToKeysNames[cidi] = map[string][]byte{key: name}
	}
	return OK, nil
}

// ERem removes the client from an ekg's set of things it's watching, and
// removes the ekg's information from the client's set of ekgs its hooked up to
func ERem(c stypes.Client, cmd *types.ClientCommand) (interface{}, error) {
	key := string(cmd.StorageKey)
	cidi := c.ClientId().Uint64()
	ekgLock.Lock()
	defer ekgLock.Unlock()

	// If the client isn't on the ekg, don't bother
	clientIdsM, ok := ekgKeyToClientIdsNames[key]
	if !ok {
		return OK, nil
	}

	if len(clientIdsM) == 1 {
		delete(ekgKeyToClientIdsNames, key)
	} else {
		delete(clientIdsM, cidi)
	}

	keysM := ekgClientIdToKeysNames[cidi]
	if len(keysM) == 1 {
		delete(ekgClientIdToKeysNames, cidi)
	} else {
		delete(keysM, key)
	}

	return OK, nil
}

// EMembers returns the list of ids being monitored by an ekg
func EMembers(c stypes.Client, cmd *types.ClientCommand) (interface{}, error) {
	key := string(cmd.StorageKey)
	ekgLock.RLock()
	defer ekgLock.RUnlock()

	clientIdsM, ok := ekgKeyToClientIdsNames[key]
	if !ok {
		return [][]byte{}, nil
	}

	names := make([][]byte, 0, len(clientIdsM))
	for _, name := range clientIdsM {
		names = append(names, name)
	}
	return names, nil
}

// ECard returns the number of client/id combinations being monitored
func ECard(c stypes.Client, cmd *types.ClientCommand) (interface{}, error) {
	key := string(cmd.StorageKey)
	ekgLock.RLock()
	defer ekgLock.RUnlock()

	clientIdsM, ok := ekgKeyToClientIdsNames[key]
	if !ok {
		return 0, nil
	}
	return len(clientIdsM), nil
}

// EkgsForClient returns a list of all the ekgs a particular client is hooked up
// to, and all the ids the client is associated with for those ekgs
func EkgsForClient(c stypes.Client) ([][]byte, [][]byte, error) {
	cidi := c.ClientId().Uint64()
	ekgLock.RLock()
	defer ekgLock.RUnlock()

	keysM, ok := ekgClientIdToKeysNames[cidi]
	if !ok {
		empty := [][]byte{}
		return empty, empty, nil
	}

	ekgs := make([][]byte, 0, len(keysM))
	ids := make([][]byte, 0, len(keysM))
	for key, id := range keysM {
		ekgs = append(ekgs, []byte(key))
		ids = append(ids, id)
	}
	return ekgs, ids, nil
}

// CleanClientEkgs takes in a client and cleans up all of the given ekgs for it,
// and the set which keeps track of those ekgs.
func CleanClientEkgs(c stypes.Client) error {
	ekgs, _, err := EkgsForClient(c)
	if err != nil {
		return err
	}
	return CleanClientEkgsShort(ekgs, c)
}

// Shortcut for CleanClientEkgs is we've already called EkgsForClient before and
// we simply want to pass that result in and not call it again. Note that this
// function deletes all record of ekgs for the given client, so the ekgs passed
// in must comprise ALL the ekgs the client is hooked up to
func CleanClientEkgsShort(ekgs [][]byte, c stypes.Client) error {
	cidi := c.ClientId().Uint64()
	ekgLock.Lock()
	defer ekgLock.Unlock()

	for _, keyb := range ekgs {
		delete(ekgKeyToClientIdsNames[string(keyb)], cidi)
	}
	delete(ekgClientIdToKeysNames, cidi)

	return nil
}
