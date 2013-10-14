package config

import (
	"fmt"
	"github.com/mediocregopher/flagconfig/src/flagconfig"
	"strings"
)

// The address hyrax should listen on for other hyrax nodes to connect
var MeshListenAddr string

// The address to advertise to other hyrax nodes that we're listening for new
// mesh connections on
var MeshAdvertiseAddr string

// Whether or not this is the first node active in the pool. If it is then
// StorageAddr and InitSecrets will be populated as well
var FirstNode bool

// The address this hyrax node, and others should use to connect to this one's
// backend storage (redis).
var StorageAddr string

// Initial secrets to load in if this is the first-node
var InitSecrets [][]byte

type ListenType int
const (
	LTYPE_TCP ListenType = iota
)

type ListenFormat int
const (
	LFORMAT_JSON ListenFormat = iota
)

// ListenAddr is a structure containing all the information needed to create a
// listener
type ListenAddr struct {

	// The type of the listen address. At the moment the only option is tcp
	Type ListenType

	// The actual address to listen for client connections on
	Addr string

	// The format to expect data to come in as. At the moment the only option is
	// json
	Format ListenFormat
}

// The list of currently active ListenAddrs
var ListenAddrs []ListenAddr

func Load() error {
	fc := flagconfig.New("hyrax")
	fc.RequiredStrParam(
		"mesh-listen-addr",
		"The address hyrax should listen on for other hyrax nodes to connect (exemple: \":9379\")",
	)
	fc.RequiredStrParam(
		"mesh-advertise-addr",
		"The address to advertise to other hyrax nodes that we're listening for new mesh connections on (example: \"127.0.0.1:9379\")",
	)
	fc.FlagParam(
		"first-node",
		"Set if this is the first node active in the pool. If it is then the init-secrets parameter is required as well",
		false,
	)
	fc.StrParams(
		"init-secrets",
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
	las := make([]ListenAddr,len(lasRaw))
	for i := range lasRaw {
		la, err := parseListenAddr(lasRaw[i])
		if err != nil {
			return err
		}
		las[i] = *la
	}

	fn := fc.GetFlag("first-node")
	isRaw := fc.GetStrs("init-secrets")
	if fn && len(isRaw) == 0 {
		return fmt.Errorf("first-node set but no init-secrets specified")
	}

	is := make([][]byte, len(isRaw))
	for i := range isRaw {
		is[i] = []byte(isRaw[i])
	}

	MeshListenAddr = fc.GetStr("mesh-listen-addr")
	MeshAdvertiseAddr = fc.GetStr("mesh-advertise-addr")
	FirstNode = fn
	InitSecrets = is
	StorageAddr = fc.GetStr("storage-addr")
	ListenAddrs = las
	return nil
}

func parseListenAddr(param string) (*ListenAddr, error) {
	pieces := strings.SplitN(param,"::",3)
	la := ListenAddr{}

	switch strings.ToLower(pieces[0]) {
	case "tcp": la.Type = LTYPE_TCP
	default: return nil, fmt.Errorf("Unknown listen type: %s", pieces[0])
	}

	la.Addr = pieces[1]

	switch strings.ToLower(pieces[2]) {
	case "json": la.Format = LFORMAT_JSON
	default: return nil, fmt.Errorf("Unknown listen format: %s", pieces[2])
	}

	return &la, nil
}
