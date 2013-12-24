package builtin

import (
	storage "github.com/mediocregopher/hyrax/server/storage-router"
	stypes "github.com/mediocregopher/hyrax/server/types"
	"github.com/mediocregopher/hyrax/types"
)

var monns = []byte("mon")

//MAdd adds the client's id to the set of clients that are monitoring the key
//(so it can receive alerts) and adds the key to the set of keys that the client
//is monitoring (so it can clean up)
func MAdd(cid stypes.ClientId, cmd *types.ClientCommand) (interface{}, error) {
	key := cmd.StorageKey
	cidb := cid.Bytes()
	monKey := storage.KeyMaker.Namespace(monns, key)
	clientMonsKey := storage.KeyMaker.ClientNamespace(monns, cidb)

	clientAdd := storage.CommandFactory.GenericSetAdd(clientMonsKey, key)
	if _, err := storage.DirectedCmd(thisnode, clientAdd); err != nil {
		return nil, err
	}

	monAdd := storage.CommandFactory.GenericSetAdd(monKey, cidb)
	return storage.DirectedCmd(thisnode, monAdd)
}

// MRem removes the client's id from the set of clients that are monitoring the
// key, and removes the key from the set of keys that the client is monitoring
func MRem(cid stypes.ClientId, cmd *types.ClientCommand) (interface{}, error) {
	key := cmd.StorageKey
	cidb := cid.Bytes()
	monKey := storage.KeyMaker.Namespace(monns, key)
	clientMonsKey := storage.KeyMaker.ClientNamespace(monns, cidb)

	monRem := storage.CommandFactory.GenericSetRem(monKey, cidb)
	r, err := storage.DirectedCmd(thisnode, monRem)
	if err != nil {
		return nil, err
	}

	clientRem := storage.CommandFactory.GenericSetRem(clientMonsKey, key)
	_, err = storage.DirectedCmd(thisnode, clientRem)
	return r, err
}

// CleanMons takes in a client id and cleans up all of its monitors, and the set
// which keeps track of those monitors
func CleanMons(cid stypes.ClientId) error {
	cidb := cid.Bytes()
	clientMonsKey := storage.KeyMaker.ClientNamespace(monns, cidb)
	monlistCmd := storage.CommandFactory.GenericSetMembers(clientMonsKey)
	r, err := storage.DirectedCmd(thisnode, monlistCmd)
	if err != nil {
		return err
	}

	mons := r.([][]byte)
	for i := range mons {
		key := mons[i]
		monKey := storage.KeyMaker.Namespace(monns, key)
		cleanKeyCmd := storage.CommandFactory.GenericSetRem(monKey, cidb)
		if _, err = storage.DirectedCmd(thisnode, cleanKeyCmd); err != nil {
			return err
		}
	}

	delClientMonCmd := storage.CommandFactory.GenericSetDel(clientMonsKey)
	_, err = storage.DirectedCmd(thisnode, delClientMonCmd)
	return err
}

// ClientsForMon takes in a key and returns all the client ids on this node that
// are mon'ing that key
func ClientsForMon(key []byte) ([]stypes.ClientId, error) {
	monKey := storage.KeyMaker.Namespace(monns, key)
	monsCmd := storage.CommandFactory.GenericSetMembers(monKey)
	r, err := storage.DirectedCmd(thisnode, monsCmd)
	if err != nil {
		return nil, err
	}

	ids := r.([][]byte)
	cids := make([]stypes.ClientId, len(ids))
	for i := range ids {
		cids[i], err = stypes.ClientIdFromBytes(ids[i])
		if err != nil {
			return nil, err
		}
	}
	return cids, nil
}
