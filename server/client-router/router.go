package router

import (
	"errors"
	stypes "github.com/mediocregopher/hyrax/server/types"
	"sync"
)

var router = map[uint64]stypes.Client{}
var routerL = sync.RWMutex{}

// Add takes in a client and adds it to the global set of clients. Returns an
// error if the client's id was already set in the router
func Add(c stypes.Client) error {
	routerL.Lock()
	defer routerL.Unlock()

	id := c.ClientId().Uint64()
	if _, ok := router[id]; ok {
		return errors.New("id already set in router")
	}

	router[id] = c
	return nil
}

// RemById takes in a client id and removes the client in the router associated
// with it, assuming one was associated
func RemById(cid stypes.ClientId) {
	routerL.Lock()
	defer routerL.Unlock()

	delete(router, cid.Uint64())
}

// RemByClient takes in a client and removes it from the router, assuming it was
// in the router in the first place
func RemByClient(c stypes.Client) {
	routerL.Lock()
	defer routerL.Unlock()

	delete(router, c.ClientId().Uint64())
}

// Get attempts to retrieve a client based on its ClientId from the router
func Get(cid stypes.ClientId) (stypes.Client, bool) {
	routerL.RLock()
	defer routerL.RUnlock()
	c, ok := router[cid.Uint64()]
	return c, ok
}
