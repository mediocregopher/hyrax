package config

import (
	"fmt"
	"log"
	"github.com/mediocregopher/flagconfig"
	"strings"
)

// Address to listen for incoming push events on
var IncomingListenAddr string

// The set of other node's IncomingListenAddr this node should connect to
var PushTo []string

// Address to listen for connections which will receive outoing push events
var OutgoingListenAddr string

// The set of other nodes' OutoingListenAddrs this node should connect to
var PullFrom []string

// The address this hyrax node, and others should use to connect to this one's
// backend storage (redis).
var StorageAddr string

// Initial secrets to load in if this is the first-node
var InitSecrets [][]byte

// ListenAddr is a structure containing all the information needed to create a
// listener
type ListenAddr struct {

	// The type of the listen address. At the moment the only option is tcp
	Type string

	// The format to expect data to come in as. At the moment the only option is
	// json
	Format string

	// The actual address to listen for client connections on
	Addr string
}

// The list of currently active ListenAddrs
var ListenAddrs []ListenAddr

func init() {
	if err := Load(); err != nil {
		log.Fatal(err)
	}
}

func Load() error {
	fc := flagconfig.New("hyrax")
	fc.StrParam(
		"incoming-listen-addr",
		"The address hyrax should listen on for connections which will send push events",
		":9379",
	)
	fc.StrParams(
		"push-to",
		"Address of another node's incoming-listen-addr that this node will forward push events to",
	)
	fc.StrParam(
		"outgoing-listen-addr",
		"The address hyrax should listen on for connections which will have push events sent to them",
		":9479",
	)
	fc.StrParams(
		"pull-from",
		"Address of another node's outgoing-listen-addr that this node will receive push events from. Can be specified multiple times",
	)
	fc.StrParams(
		"init-secret",
		"A global secret key as a string. Can be specified multiple times if this is a first-node",
	)
	fc.StrParam(
		"storage-addr",
		"The address this hyrax node, and others should use to connect to this one's backend storage (redis).",
		"127.0.0.1:6379",
	)
	fc.StrParams(
		"listen-addr",
		"The type, address, and format to listen for client connections on, separated by a \"::\". At the moment the only type is tcp, the only format is json. Can be specified multiple times",
		"tcp::json:::2379",
	)
	if err := fc.Parse(); err != nil {
		return err
	}

	lasRaw := fc.GetStrs("listen-addr")
	las := make([]ListenAddr, len(lasRaw))
	for i := range lasRaw {
		la, err := parseListenAddr(lasRaw[i])
		if err != nil {
			return err
		}
		las[i] = *la
	}

	fn := fc.GetFlag("first-node")
	isRaw := fc.GetStrs("init-secret")
	if fn && len(isRaw) == 0 {
		return fmt.Errorf("first-node set but no init-secrets specified")
	}

	is := make([][]byte, len(isRaw))
	for i := range isRaw {
		is[i] = []byte(isRaw[i])
	}

	IncomingListenAddr = fc.GetStr("incoming-listen-addr")
	PushTo = fc.GetStrs("push-to")
	OutgoingListenAddr = fc.GetStr("outgoing-listen-addr")
	PullFrom = fc.GetStrs("pull-from")
	InitSecrets = is
	StorageAddr = fc.GetStr("storage-addr")
	ListenAddrs = las
	return nil
}

func parseListenAddr(param string) (*ListenAddr, error) {
	pieces := strings.SplitN(param, "::", 3)
	la := ListenAddr{
		Type:   pieces[0],
		Format: pieces[1],
		Addr:   pieces[2],
	}

	return &la, nil
}
