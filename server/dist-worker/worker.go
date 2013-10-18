package distworker

import (
	"github.com/mediocregopher/hyrax/server/dist"
	"github.com/mediocregopher/hyrax/server/builtin"
	crouter "github.com/mediocregopher/hyrax/server/client-router"
	storage "github.com/mediocregopher/hyrax/server/storage-router"
)

func init() {
	go spin()
}

func spin() {
	for {
		select {
		case sbs := <- dist.StorageBucketSetCommands:
			// TODO Do something with error
			storage.SetBucket(sbs.Bucket, sbs.ConnType, sbs.Addr, sbs.Extra...)

		case cmd := <- dist.ClientCommands:
			cmd.Secret = nil
			// TODO Do something with the error
			cids, _ := builtin.ClientsForMon(cmd.StorageKey)
			for i := range cids {
				if c, ok := crouter.Get(cids[i]); ok {
					c.CommandPushChannel() <- cmd
				} else {
					// TODO error because client with id wasn't found
				}
			}
		}
	}
}
