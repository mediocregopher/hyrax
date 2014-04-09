package pubsub

import (
	"github.com/grooveshark/golib/gslog"
	"sync"
	"time"

	stypes "github.com/mediocregopher/hyrax/server/types"
	"github.com/mediocregopher/hyrax/types"
)

const PUB_CHUNK_SIZE = 500

// A system wherein clients can subscribe to channels and others can publish to
// those channels. Each PubSub instance is a totally separate system, they do
// not overlap in anyway.
type PubSub struct {
	subClients map[string]map[stypes.Client]bool
	clientSubs map[stypes.Client]map[string]bool
	subChs     map[string]chan *types.Action
	subLock    sync.RWMutex
}

// Returns a new PubSub system
func New() *PubSub {
	return &PubSub{
		subClients: map[string]map[stypes.Client]bool{},
		clientSubs: map[stypes.Client]map[string]bool{},
		subChs:     map[string]chan *types.Action{},
	}
}

func (ps *PubSub) subSpin(sub string) {
	ps.subLock.RLock()
	subCh, ok1 := ps.subChs[sub]
	clients, ok2 := ps.subClients[sub]
	ps.subLock.RUnlock()

	if !ok1 || !ok2 {
		gslog.Errorf("Missing data for sub %s in pubsub", sub)
		return
	}

	for cmd := range subCh {
		ps.subLock.RLock()
		if len(clients) < PUB_CHUNK_SIZE {
			for client := range clients {
				pubToClient(client, cmd, sub)
			}
		} else {
			clientCh := make(chan stypes.Client)
			go chunker(clientCh, cmd, sub)
			for client := range clients {
				clientCh <- client
			}
			close(clientCh)
		}
		ps.subLock.RUnlock()
	}
}

func chunker(clientCh <-chan stypes.Client, cmd *types.Action, sub string) {
outer:
	for i := 0; ; {
		chunkCh := make(chan stypes.Client, PUB_CHUNK_SIZE)
		go chunkPubber(chunkCh, cmd, sub)
		for client := range clientCh {
			chunkCh <- client
			i++
			if i%PUB_CHUNK_SIZE == 0 {
				close(chunkCh)
				continue outer
			}
		}
		close(chunkCh)
		return
	}
}

func chunkPubber(chunkCh <-chan stypes.Client, cmd *types.Action, sub string) {
	for c := range chunkCh {
		pubToClient(c, cmd, sub)
	}
}

func pubToClient(c stypes.Client, cmd *types.Action, sub string) {
	select {
	case c.PushCh() <- cmd:
	case <-time.After(10 * time.Second):
		gslog.Warnf("Timeout pubbing to %p for sub %s", c, sub)
	}
}

// Subscribes a client so that they will receive push messages on the given
// subscriptions
func (ps *PubSub) Subscribe(cl stypes.Client, subs ...string) error {
	ps.subLock.Lock()
	defer ps.subLock.Unlock()

	for _, sub := range subs {
		sc, ok := ps.subClients[sub]
		if ok {
			sc[cl] = true
		} else {
			ps.subClients[sub] = map[stypes.Client]bool{cl: true}
			subCh := make(chan *types.Action)
			ps.subChs[sub] = subCh
			go ps.subSpin(sub)
		}
	}

	cs, ok := ps.clientSubs[cl]
	if !ok {
		cs = map[string]bool{}
		ps.clientSubs[cl] = cs
	}
	for _, sub := range subs {
		cs[sub] = true
	}
	return nil
}

// Unsubscribes a client so from the given subscriptions
func (ps *PubSub) Unsubscribe(cl stypes.Client, subs ...string) error {
	ps.subLock.Lock()
	defer ps.subLock.Unlock()

	// If the client doesn't appear in clientSubs then it has no subs
	cs, ok := ps.clientSubs[cl]
	if !ok {
		return nil
	}

	for _, sub := range subs {
		delete(cs, sub)
		sc, ok := ps.subClients[sub]
		if !ok {
			gslog.Errorf("No subClients for sub %s", sub)
			continue
		}
		delete(sc, cl)
		if len(sc) == 0 {
			delete(ps.subClients, sub)
			close(ps.subChs[sub])
			delete(ps.subChs, sub)
		}
	}

	if len(cs) == 0 {
		delete(ps.clientSubs, cl)
	}

	return nil
}

// Returns all the subscriptions a client is subscribed to
func (ps *PubSub) GetSubscriptions(cl stypes.Client) ([]string, error) {
	ps.subLock.RLock()
	defer ps.subLock.RUnlock()

	cs, ok := ps.clientSubs[cl]
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
func (ps *PubSub) UnsubscribeAll(cs stypes.Client) error {
	subs, err := ps.GetSubscriptions(cs)
	if err != nil {
		return err
	}
	return ps.Unsubscribe(cs, subs...)
}

// Publishes the given command to all clients subscribed to the given
// subscriptions
func (ps *PubSub) Publish(a *types.Action, subs ...string) error {
	ps.subLock.RLock()
	defer ps.subLock.RUnlock()

	for _, sub := range subs {
		subCh, ok := ps.subChs[sub]
		if !ok {
			continue
		}
		subCh <- a
	}

	return nil
}
