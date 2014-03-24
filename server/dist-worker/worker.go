package distworker

import (
	"github.com/mediocregopher/hyrax/server/core"
	crouter "github.com/mediocregopher/hyrax/server/client-router"
	"github.com/mediocregopher/hyrax/server/dist"
	storage "github.com/mediocregopher/hyrax/server/storage-router"
)

func Init() {
	go spin()
}

func spin() {
	for {
		select {
		case sbs := <-dist.StorageBucketSetCommands:
			// TODO Do something with error
			storage.SetBucket(sbs.Bucket, sbs.ConnType, sbs.Addr, sbs.Extra...)

		case cmd := <-dist.ClientCommands:
			cmd.Secret = nil
			// TODO Do something with the error
			cids, _ := core.ClientIdsForMon(cmd.StorageKey)
			for i := range cids {
				if c, ok := crouter.Get(cids[i]); ok {
					c.CommandPushCh() <- cmd
				} else {
					// TODO error because client with id wasn't found
				}
			}
		}
	}
}
