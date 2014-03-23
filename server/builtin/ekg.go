package builtin

import (
	"sync"

	stypes "github.com/mediocregopher/hyrax/server/types"
	"github.com/mediocregopher/hyrax/types"
)

var ekgns = []byte("ekg")

// A mapping of ekgs to ClientIds and their names
var ekgKeyToClientIdsNames = map[string]map[uint64][]byte{}

// A mapping of ClientIds to the ekgs the client is on and the names it's using
// for each ekg
var ekgClientIdToKeysNames = map[uint64]map[string][]byte{}

// Lock which coordinates access to the mappings
var ekgLock sync.RWMutex

// EAdd adds the client's id (actual and given) to an ekg's set of things it's
// watching, and adds the ekg's information to the client's set of ekgs its
// hooked up to
func EAdd(cid stypes.ClientId, cmd *types.ClientCommand) (interface{}, error) {
	key := string(cmd.StorageKey)
	cidi := cid.Uint64()
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
	return []byte("OK"), nil
}

// ERem removes the client's id from an ekg's set of things it's watching, and
// removes the ekg's information from the client's set of ekgs its hooked up to
func ERem(cid stypes.ClientId, cmd *types.ClientCommand) (interface{}, error) {
	key := string(cmd.StorageKey)
	cidi := cid.Uint64()
	ekgLock.Lock()
	defer ekgLock.Unlock()

	// If the client isn't on the ekg, don't bother
	clientIdsM, ok := ekgKeyToClientIdsNames[key]
	if !ok {
		return []byte("OK"), nil
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

	return []byte("OK"), nil
}

// EMembers returns the list of ids being monitored by an ekg
func EMembers(
	cid stypes.ClientId,
	cmd *types.ClientCommand) (interface{}, error) {

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
func ECard(
	cid stypes.ClientId,
	cmd *types.ClientCommand) (interface{}, error) {

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
func EkgsForClient(cid stypes.ClientId) ([][]byte, [][]byte, error) {
	cidi := cid.Uint64()
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

// CleanClientEkgs takes in a client id and cleans up all of the given ekgs for
// it, and the set which keeps track of those ekgs.
func CleanClientEkgs(cid stypes.ClientId) error {
	ekgs, _, err := EkgsForClient(cid)
	if err != nil {
		return err
	}
	return CleanClientEkgsShort(ekgs, cid)
}

// Shortcut for CleanClientEkgs is we've already called EkgsForClient before and
// we simply want to pass that result in and not call it again. Note that this
// function deletes all record of ekgs for the given client id, so the ekgs
// passed in must comprise ALL the ekgs the client is hooked up to
func CleanClientEkgsShort(ekgs [][]byte, cid stypes.ClientId) error {
	cidi := cid.Uint64()
	ekgLock.Lock()
	defer ekgLock.Unlock()

	for _, keyb := range ekgs {
		delete(ekgKeyToClientIdsNames[string(keyb)], cidi)
	}
	delete(ekgClientIdToKeysNames, cidi)

	return nil
}
