package router

import (
	"time"
	"sync"

	stypes "github.com/mediocregopher/hyrax/server/types"
	"github.com/mediocregopher/hyrax/types"
)

var subClients = map[string]map[stypes.Client]bool{}
var clientSubs = map[stypes.Client]map[string]bool{}
var subChs = map[string]chan *types.ClientCommand{}
var subLock sync.RWMutex

func subSpin(sub string) {
	subLock.RLock()
	subCh, ok1 := subChs[sub]
	clients, ok2 := subClients[sub]
	subLock.RUnlock()

	if !ok1 || !ok2 {
		// TODO log error
		return
	}

	for cmd := range subCh {
		subLock.RLock()
		for client := range clients {
			select {
			case client.CommandPushCh() <- cmd:
			case <-time.After(10 * time.Second):
				// TODO log error
			}
		}
		subLock.RUnlock()
	}
}

// Subscribes a client so that they will receive push messages on the given
// subscriptions
func Subscribe(cl stypes.Client, subs ...string) error {
	subLock.Lock()
	defer subLock.Unlock()

	for _, sub := range subs {
		sc, ok := subClients[sub]
		if ok {
			sc[cl] = true
		} else {
			subClients[sub] = map[stypes.Client]bool{cl: true}
			subCh := make(chan *types.ClientCommand)
			subChs[sub] = subCh
			go subSpin(sub)
		}
	}

	cs, ok := clientSubs[cl]
	if !ok {
		cs = map[string]bool{}
		clientSubs[cl] = cs
	}
	for _, sub := range subs {
		cs[sub] = true
	}
	return nil
}

// Unsubscribes a client so from the given subscriptions
func Unsubscribe(cl stypes.Client, subs ...string) error {
	subLock.Lock()
	defer subLock.Unlock()

	// If the client doesn't appear in clientSubs then it has no subs
	cs, ok := clientSubs[cl]
	if !ok {
		return nil
	}

	for _, sub := range subs {
		delete(cs, sub)
		sc, ok := subClients[sub]
		if !ok {
			// TODO log error
			continue
		}
		delete(sc, cl)
		if len(sc) == 0 {
			delete(subClients, sub)
			close(subChs[sub])
			delete(subChs, sub)
		}
	}

	if len(cs) == 0 {
		delete(clientSubs, cl)
	}

	return nil
}

// Returns all the subscriptions a client is subscribed to
func GetSubscriptions(cl stypes.Client) ([]string, error) {
	subLock.RLock()
	defer subLock.RUnlock()

	cs, ok := clientSubs[cl]
	if !ok {
		return []string{}, nil
	}

	ret := make([]string, 0, len(cs))
	for sub := range cs {
		ret = append(ret, sub)
	}
	return ret, nil
}

// Unsubscribes a given client from any subscriptions it's subscribed to
func UnsubscribeAll(cs stypes.Client) error {
	subs, err := GetSubscriptions(cs)
	if err != nil {
		return err
	}
	return Unsubscribe(cs, subs...)
}

// Publishes the given command to all clients subscribed to the given
// subscriptions
func Publish(cmd *types.ClientCommand, subs ...string) error {
	subLock.RLock()
	defer subLock.RUnlock()

	for _, sub := range subs {
		subCh, ok := subChs[sub]
		if !ok {
			continue
		}
		subCh <- cmd
	}

	return nil
}
