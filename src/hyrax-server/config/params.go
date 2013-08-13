package config

//param Type
const (
	STRING = iota
	INT
)

type param struct {
	Description string
	Type        int
	Default     interface{}
}

// The list of parameters that hyrax will use during runtime, and all of their
// metadata
var params = map[string]param{
	"port": param{Description: "The tcp port to listen for new connections on",
		Type:    INT,
		Default: 3400},

	"redis-addr": param{Description: "The hostname:port (or unix sock location) to connect to redis on",
		Type:    STRING,
		Default: "localhost:6379"},

	"initial-secret-keys": param{Description: "The list of colon-separated secret-keys to start with. Can be changed on-the-fly later. Empty string for no keys at all",
		Type:    STRING,
		Default: "change:me"},
}
