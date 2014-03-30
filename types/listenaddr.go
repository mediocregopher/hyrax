package types

import (
	"strings"
)

// ListenAddr is a structure containing all the information needed to create or
// connect to a hyrax endpoint
type ListenAddr struct {

	// The type of the listen address. At the moment the only option is tcp
	Type string

	// The format to expect data to come in as. At the moment the only option is
	// json
	Format string

	// The actual address to listen for client connections on
	Addr string
}

func NewListenAddr(conntype, format, addr string) *ListenAddr {
	return &ListenAddr{conntype, format, addr}
}

// Takes a flat string and parses out a ListenAddr
func ListenAddrFromString(param string) (*ListenAddr, error) {
	pieces := strings.SplitN(param, "::", 3)
	la := ListenAddr{
		Type:   pieces[0],
		Format: pieces[1],
		Addr:   pieces[2],
	}

	return &la, nil
}

func (la *ListenAddr) String() string {
	parts := []string{la.Type, la.Format, la.Addr}
	return strings.Join(parts, "::")
}
