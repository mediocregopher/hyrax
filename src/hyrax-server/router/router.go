package router

import (
	"github.com/mediocregopher/hyrax/src/hyrax/types"
	"github.com/mediocregopher/hyrax/src/hyrax-server/router/storage"
	"github.com/mediocregopher/hyrax/src/hyrax-server/router/storage/command"
	"github.com/mediocregopher/hyrax/src/hyrax-server/router/bucket"
)

// SetBucket sets the given bucket index to be the connection to the given
// storage unit, whose name is given by its address
func SetBucket(bIndex int, conntype, addr string, extra ...interface{}) error {
	if err := storage.AddUnit(addr, conntype, addr, extra...); err != nil {
		return err
	}

	if err := bucket.Set(&addr, bIndex); err != nil {
		return err
	}

	// TODO: Check for orphaned storage unit connections, right now if a
	// connection is no longer in the bucket list it will still hold its
	// connection

	return nil
}

// GetBuckets returns a copy of the current bucket 
func GetBuckets() []*string {
	return bucket.Buckets()
}

// RoutedCmd takes in the key for a command and the command to perform. The
// Command will probably contain the key within it, the key passed here is used
// only for routing the command to the proper storage unit. The return from the
// command and an error are returned.
func RoutedCmd(key types.Byter, cmd command.Command) (interface{}, error) {
	b, err := bucket.KeyBucket(key)
	if err != nil {
		return nil, err
	}

	return storage.Cmd(b, cmd)
}

// DirectedCmd takes in the addr used to identify a storage unit and the command
// to perform on that unit, and returns the return from the command or an error
func DirectedCmd(addr *string, cmd command.Command) (interface{}, error) {
	return storage.Cmd(addr, cmd)
}
