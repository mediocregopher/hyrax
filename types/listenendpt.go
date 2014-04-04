package types

import (
	"strings"
)

// ListenEndpoint is a structure containing all the information needed to create
// or connect to a hyrax endpoint
type ListenEndpoint struct {

	// The type of the endpoint. At the moment the only option is tcp
	Type string

	// The format to expect data to come in as. At the moment the only option is
	// json
	Format string

	// The actual address to listen for client connections on
	Addr string
}

func NewListenEndpoint(conntype, format, addr string) *ListenEndpoint {
	return &ListenEndpoint{conntype, format, addr}
}

// Takes a flat string and parses out a ListenEndpoint
func ListenEndpointFromString(param string) (*ListenEndpoint, error) {
	pieces := strings.SplitN(param, "::", 3)
	le := ListenEndpoint{
		Type:   pieces[0],
		Format: pieces[1],
		Addr:   pieces[2],
	}

	return &le, nil
}

func (le *ListenEndpoint) String() string {
	parts := []string{le.Type, le.Format, le.Addr}
	return strings.Join(parts, "::")
}
