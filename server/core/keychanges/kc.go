package keychanges

import (
	"github.com/mediocregopher/hyrax/server/pubsub"
	stypes "github.com/mediocregopher/hyrax/server/types"
	"github.com/mediocregopher/hyrax/types"
)

// The channel to use in the PubSub objects when we're only using the PubSub as
// a giant single channel (like in global and local)
const single = ""

var global = pubsub.New()
var local = pubsub.New()
var mon = pubsub.New()

// Subscribes a client to global key change events. These are events which are
// being broadcast out to every node in the cluster.
func SubscribeGlobal(cl stypes.Client) error {
	return global.Subscribe(cl, single)
}

// Unsubscribes a client from receiving global key change events, if it was
// receiving any at all
func UnsubscribeGlobal(cl stypes.Client) error {
	return global.Unsubscribe(cl, single)
}

// Publishes a key change globally, both to those subscribed to global key
// changes and those subscribed (mon'd) to the actual key being changed
func PubGlobal(cc *types.ClientCommand) error {
	if err := global.Publish(cc, single); err != nil {
		return err
	}

	return mon.Publish(cc, cc.StorageKey)
}

// Subscribes a client to local key change events, which are events that
// originated on this server.
func SubscribeLocal(cl stypes.Client) error {
	return local.Subscribe(cl, single)
}

// Unsubscribes a client from receiving local key change events, if it was
// receiving any at all
func UnsubscribeLocal(cl stypes.Client) error {
	return local.Unsubscribe(cl, single)
}

// Publishes a key change to those subscribed to local key change events
func PubLocal(cc *types.ClientCommand) error {
	return local.Publish(cc, single)
}

// Subscribes a client to receive keychange events about a particular key
func Mon(cl stypes.Client, keys ...string) error {
	keysStr := make([]string, len(keys))
	for i := range keys {
		keysStr[i] = keys[i]
	}
	return mon.Subscribe(cl, keysStr...)
}

// Unsubscribes a client from particular keys, if it was subscribed at all
func Unmon(cl stypes.Client, keys ...string) error {
	keysStr := make([]string, len(keys))
	for i := range keys {
		keysStr[i] = keys[i]
	}
	return mon.Unsubscribe(cl, keysStr...)
}

// Unsubscribes a client from any key change events it might be receiving
func UnsubscribeAll(cl stypes.Client) error {
	if err := global.Unsubscribe(cl, single); err != nil {
		return err
	}

	if err := local.Unsubscribe(cl, single); err != nil {
		return err
	}

	if err := mon.UnsubscribeAll(cl); err != nil {
		return err
	}
	return nil
}
