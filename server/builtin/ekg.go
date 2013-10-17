package builtin

import (
	"github.com/mediocregopher/hyrax/server/config"
	storage "github.com/mediocregopher/hyrax/server/storage-router"
	"github.com/mediocregopher/hyrax/types"
	ctypes "github.com/mediocregopher/hyrax/types/client"
	stypes "github.com/mediocregopher/hyrax/server/types"
)

var ekgns = types.SimpleByter([]byte("ekg"))

// Some shortcuts
var keyMaker = storage.KeyMaker
var cmdFactory = storage.CommandFactory

// EAdd adds the client's id (actual and given) to an ekg's set of things it's
// watching, and adds the ekg's information to the client's set of ekgs its
// hooked up to
func EAdd(cid stypes.ClientId, cmd *ctypes.ClientCommand) (interface{}, error) {
	key := cmd.StorageKey	
	id := types.NewByter(cmd.Id)
	ekgKey := keyMaker.Namespace(ekgns, key)
	clientEkgsKey := keyMaker.ClientNamespace(ekgns, cid)
	thisnode := &config.StorageAddr
	
	clAdd := cmdFactory.GenericSetAdd(clientEkgsKey, key)
	if _, err := storage.DirectedCmd(thisnode, clAdd); err != nil {
		return nil, err
	}

	addCmd := storage.CommandFactory.KeyValueSetAdd(ekgKey, cid, id)
	return storage.RoutedCmd(key, addCmd)
}

// ERem removes the client's id from an ekg's set of things it's watching, and
// removes the ekg's information from the client's set of ekgs its hooked up to
func ERem(cid stypes.ClientId, cmd *ctypes.ClientCommand) (interface{}, error) {
	key := cmd.StorageKey	
	ekgKey := keyMaker.Namespace(ekgns, key)
	clientEkgsKey := keyMaker.Namespace(ekgns, cid)
	thisnode := &config.StorageAddr

	remCmd := cmdFactory.KeyValueSetRemByInnerKey(ekgKey, cid)
	r, err := storage.RoutedCmd(key, remCmd)
	if err != nil {
		return nil, err
	}

	clRem := cmdFactory.GenericSetRem(clientEkgsKey, key)
	if _, err := storage.DirectedCmd(thisnode, clRem); err != nil {
		return nil, err
	}

	return r, nil
}

// EMembers returns the list of ids being monitored by an ekg
func EMembers(
	cid stypes.ClientId,
	cmd *ctypes.ClientCommand) (interface{}, error) {

	key := cmd.StorageKey
	ekgKey := keyMaker.Namespace(ekgns, key)
	memsCmd := cmdFactory.KeyValueSetMemberValues(ekgKey)
	return storage.RoutedCmd(key, memsCmd)
}

// ECard returns the number of client/id combinations being monitored
func ECard(
	cid stypes.ClientId,
	cmd *ctypes.ClientCommand) (interface{}, error) {

	key := cmd.StorageKey
	ekgKey := keyMaker.Namespace(ekgns, key)
	cardCmd := cmdFactory.KeyValueSetCard(ekgKey)
	return storage.RoutedCmd(key, cardCmd)
}

// EkgsForClient returns a list of all the ekgs a particular client is hooked up
// to
func EkgsForClient(cid stypes.ClientId) ([]types.Byter, error) {
	clientEkgsKey := keyMaker.ClientNamespace(ekgns, cid)
	ekgsCmd := cmdFactory.GenericSetMembers(clientEkgsKey)
	thisnode := &config.StorageAddr
	r, err := storage.DirectedCmd(thisnode, ekgsCmd)
	if err != nil {
		return nil, err
	}

	ekgs := r.([][]byte)
	ekgsb := make([]types.Byter, len(ekgs))
	for i := range ekgs {
		ekgsb[i] = types.NewByter(ekgs[i])
	}

	return ekgsb, nil
}

// CleanClientEkgs takes in a client id and cleans up all of its ekgs, and the
// set which keeps track of those ekgs.
func CleanClientEkgs(cid stypes.ClientId) error {
	ekgs, err := EkgsForClient(cid)
	if err != nil {
		return err
	}

	for i := range ekgs {
		key := ekgs[i]
		ekgKey := keyMaker.Namespace(ekgns, key)
		remCmd := cmdFactory.KeyValueSetRemByInnerKey(ekgKey, cid)
		if _, err = storage.RoutedCmd(key, remCmd); err != nil {
			return err
		}
	}

	return nil
}
