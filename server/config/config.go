package config

import (
	"github.com/grooveshark/golib/gslog"
	"github.com/mediocregopher/flagconfig"

	"github.com/mediocregopher/hyrax/types"
)

// Information for connecting to the storage instance
var StorageInfo string

// Initial secrets to load in if this is the first-node
var InitSecrets [][]byte

// The list of endpoints this node should server
var ListenEndpoints []types.ListenEndpoint

// The list of endpoints this node will send local key change events to
var PushToEndpoints []types.ListenEndpoint

// The list of endpoints this node will pull global key change events from
var PullFromEndpoints []types.ListenEndpoint

// The endpoint to advertise to other nodes that they should connect to
var MyEndpoint types.ListenEndpoint

func init() {
	if err := Load(); err != nil {
		gslog.Fatal(err.Error())
	}
}

func Load() error {
	fc := flagconfig.New("hyrax")
	fc.StrParams(
		"init-secret",
		"A global secret key as a string. Can be specified multiple times if this is a first-node",
	)
	fc.StrParam(
		"storage-info",
		"Info needed for connecting to the datastore(s). For redis this is just the address redis is listening on",
		"127.0.0.1:6379",
	)
	fc.StrParams(
		"listen-endpoint",
		"The type, address, and format to listen for client connections on, separated by a \"::\". At the moment the only type is tcp, the only format is json. Can be specified multiple times",
		"tcp::json:::2379",
	)
	fc.StrParams(
		"push-to-endpoint",
		"The endpoint address (see listen-endpoint for format) this node will send local keychange events to. Can be specified multiple times",
	)
	fc.StrParams(
		"pull-from-endpoint",
		"The endpoint address (see listen-endpoint for format) this node will pull global keychange events from. Can be specified multiple times",
	)
	fc.StrParam(
		"my-endpoint",
		"The endpoint address (see listen-endpoint for format) this node will advertise to other nodes that they should connect to",
		"tcp::json::localhost:2379",
	)
	if err := fc.Parse(); err != nil {
		return err
	}

	isRaw := fc.GetStrs("init-secret")
	is := make([][]byte, len(isRaw))
	for i := range isRaw {
		is[i] = []byte(isRaw[i])
	}

	InitSecrets = is
	StorageInfo = fc.GetStr("storage-info")

	var err error
	if ListenEndpoints, err = endpts(fc, "listen-endpoint"); err != nil {
		return err
	}
	if PushToEndpoints, err = endpts(fc, "push-to-endpoint"); err != nil {
		return err
	}
	if PullFromEndpoints, err = endpts(fc, "pull-from-endpoint"); err != nil {
		return err
	}

	myEndpointRaw := fc.GetStr("my-endpoint")
	myEndpoint, err := types.ListenEndpointFromString(myEndpointRaw)
	if err != nil {
		return err
	}
	MyEndpoint = *myEndpoint

	return nil
}

func endpts(
	fc *flagconfig.FlagConfig, param string) ([]types.ListenEndpoint, error) {

	lesRaw := fc.GetStrs(param)
	les := make([]types.ListenEndpoint, len(lesRaw))
	for i := range lesRaw {
		le, err := types.ListenEndpointFromString(lesRaw[i])
		if err != nil {
			return nil, err
		}
		les[i] = *le
	}
	return les, nil
}
