package keychanges

import (
	"log"
	"sync"

	"github.com/mediocregopher/hyrax/types"
	stypes "github.com/mediocregopher/hyrax/server/types"
)

// All key changes go through here
var keyChangeCh = make(chan *types.ClientCommand)

// Set of all clients and their proxy channels who are listening for key change
// events
var keyChangeClients = map[uint64]chan *types.ClientCommand{}

// For synchronizing changes to keyChangeClients
var keyChangeLock = sync.RWMutex{}

// A channel on which commands which will modify keys in the data store should
// be pushed to
func Ch() chan<- *types.ClientCommand {
	return keyChangeCh
}

func init() {
	// Spawn the spinner
	go func() {
		for cc := range keyChangeCh {
			keyChangeLock.RLock()
			for _, proxyCh := range keyChangeClients {
				proxyCh <- cc
			}
			keyChangeLock.RUnlock()
		}
	}()
}

// Registers a client as wanting key change events pushed to them
func AddClient(c stypes.Client) error {
	cid := c.ClientId().Uint64()
	keyChangeLock.Lock()
	defer keyChangeLock.Unlock()
	if _, ok := keyChangeClients[cid]; ok {
		return nil
	}

	proxyCh := make(chan *types.ClientCommand)
	go proxySpin(proxyCh, c)
	keyChangeClients[cid] = proxyCh
	return nil
}

// Unregisters a client from having key change events pushed to them
func RemoveClient(c stypes.Client) error {
	cid := c.ClientId().Uint64()
	keyChangeLock.Lock()
	defer keyChangeLock.Unlock()
	if proxyCh, ok := keyChangeClients[cid]; ok {
		close(proxyCh)
		delete(keyChangeClients, cid)
	}
	return nil
}

func proxySpin(proxyCh chan *types.ClientCommand, c stypes.Client) {
	pushCh := c.CommandPushCh()
	closeCh := c.ClosingCh()
	for cc := range proxyCh {
		select {
		case pushCh<- cc:
		case <-closeCh:
			// Must be grounded by another go-routine in order prevent
			// race-condition where spinner is writing to proxyCh while
			// RemoveClient is waiting to delete it
			go ground(proxyCh)
			if err := RemoveClient(c); err != nil {
				log.Printf("Error removing key change client: %s", err)
			}
		}
	}
}

func ground(ch chan *types.ClientCommand) {
	for _ = range ch {
	}
}
