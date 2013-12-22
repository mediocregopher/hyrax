package builtin

import (
	"github.com/mediocregopher/hyrax/server/config"
	storage "github.com/mediocregopher/hyrax/server/storage-router"
	"github.com/mediocregopher/hyrax/types"
	stypes "github.com/mediocregopher/hyrax/server/types"
)

var ekgns = types.SimpleByter([]byte("ekg"))

// Some shortcuts
var keyMaker = storage.KeyMaker
var cmdFactory = storage.CommandFactory

// EAdd adds the client's id (actual and given) to an ekg's set of things it's
// watching, and adds the ekg's information to the client's set of ekgs its
// hooked up to
func EAdd(cid stypes.ClientId, cmd *types.ClientCommand) (interface{}, error) {
	key := cmd.StorageKey	
	id := types.NewByter(cmd.Id)
	ekgKey := keyMaker.Namespace(ekgns, key)
	clientEkgsKey := keyMaker.ClientNamespace(ekgns, cid)
	thisnode := &config.StorageAddr
	
	clAdd := cmdFactory.KeyValueSetAdd(clientEkgsKey, key, id)
	if _, err := storage.DirectedCmd(thisnode, clAdd); err != nil {
		return nil, err
	}

	addCmd := storage.CommandFactory.KeyValueSetAdd(ekgKey, cid, id)
	return storage.RoutedCmd(key, addCmd)
}

// ERem removes the client's id from an ekg's set of things it's watching, and
// removes the ekg's information from the client's set of ekgs its hooked up to
func ERem(cid stypes.ClientId, cmd *types.ClientCommand) (interface{}, error) {
	key := cmd.StorageKey	
	ekgKey := keyMaker.Namespace(ekgns, key)
	clientEkgsKey := keyMaker.Namespace(ekgns, cid)
	thisnode := &config.StorageAddr

	remCmd := cmdFactory.KeyValueSetRemByInnerKey(ekgKey, cid)
	r, err := storage.RoutedCmd(key, remCmd)
	if err != nil {
		return nil, err
	}

	clRem := cmdFactory.KeyValueSetRemByInnerKey(clientEkgsKey, key)
	if _, err := storage.DirectedCmd(thisnode, clRem); err != nil {
		return nil, err
	}

	return r, nil
}

// EMembers returns the list of ids being monitored by an ekg
func EMembers(
	cid stypes.ClientId,
	cmd *types.ClientCommand) (interface{}, error) {

	key := cmd.StorageKey
	ekgKey := keyMaker.Namespace(ekgns, key)
	memsCmd := cmdFactory.KeyValueSetMemberValues(ekgKey)
	return storage.RoutedCmd(key, memsCmd)
}

// ECard returns the number of client/id combinations being monitored
func ECard(
	cid stypes.ClientId,
	cmd *types.ClientCommand) (interface{}, error) {

	key := cmd.StorageKey
	ekgKey := keyMaker.Namespace(ekgns, key)
	cardCmd := cmdFactory.KeyValueSetCard(ekgKey)
	return storage.RoutedCmd(key, cardCmd)
}

// EkgsForClient returns a list of all the ekgs a particular client is hooked up
// to, and all the ids the client is associated with for those ekgs
func EkgsForClient(cid stypes.ClientId) ([]types.Byter, [][]byte, error) {
	clientEkgsKey := keyMaker.ClientNamespace(ekgns, cid)
	ekgsCmd := cmdFactory.KeyValueSetMembers(clientEkgsKey)
	thisnode := &config.StorageAddr
	r, err := storage.DirectedCmd(thisnode, ekgsCmd)
	if err != nil {
		return nil, nil, err
	}

	rall := r.([][]byte)
	ekgs := make([]types.Byter, len(rall)/2)
	ids := make([][]byte, len(rall)/2)
	for i:=0; i<len(rall); i += 2 {
		ekgs[i/2] = types.NewByter(rall[i])
		ids[i/2] = rall[i+1]
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
func CleanClientEkgsShort(ekgs []types.Byter, cid stypes.ClientId) error {
	for i := range ekgs {
		key := ekgs[i]
		ekgKey := keyMaker.Namespace(ekgns, key)
		remCmd := cmdFactory.KeyValueSetRemByInnerKey(ekgKey, cid)
		if _, err := storage.RoutedCmd(key, remCmd); err != nil {
			return err
		}
	}

	clientEkgsKey := keyMaker.ClientNamespace(ekgns, cid)
	clRemCmd := cmdFactory.KeyValueSetDel(clientEkgsKey)
	thisnode := &config.StorageAddr
	if _, err := storage.DirectedCmd(thisnode, clRemCmd); err != nil {
		return err
	}

	return nil
}
