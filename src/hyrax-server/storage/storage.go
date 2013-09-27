package storage

import (
	"time"
	"errors"
	"github.com/mediocregopher/hyrax/src/hyrax-server/storage/redis"
	"github.com/mediocregopher/hyrax/src/hyrax-server/storage/unit"
	sucmd "github.com/mediocregopher/hyrax/src/hyrax-server/storage/command"
)

type LocatorFunc func([]*unit.StorageUnit) *unit.StorageUnit

type commandWrap struct {
	unitname *string
	cmdb *sucmd.CommandBundle
}

type storageManager struct {
	units map[string]*unit.StorageUnit
	addUnitCh chan *unit.StorageUnit
	remUnitCh chan string
	cmdCh chan *commandWrap
}

var sm = storageManager{
	units: map[string]*unit.StorageUnit{},
	addUnitCh: make(chan *unit.StorageUnit),
	remUnitCh: make(chan string),
	cmdCh: make(chan *commandWrap),
}

var NewStorageUnitConn func() unit.StorageUnitConn
var CommandFactory sucmd.CommandFactory
var NewTransaction func(...sucmd.Command) sucmd.Command

// Init starts up the storage manager and prepares various storage sepecific
// units for use by the outside world
func Init() {
	NewStorageUnitConn = redis.New
	CommandFactory = sucmd.CommandFactory(&redis.RedisCommandFactory{})
	NewTransaction = redis.NewRedisTransaction
	go sm.spin()
}

func (sm *storageManager) spin() {
	for {
		select {
		//case u := <-sm.addUnitCh:
		//case r := <-sm.remUnitCh:
		case c := <-sm.cmdCh:
			if unit, ok := sm.units[*c.unitname]; ok {
				//We assume our StorageUnits won't block on this call. If any do
				//we have big problems anyway
				unit.Cmd(c.cmdb)
			} else {
				err := errors.New("Unknown storage unit: "+*c.unitname)
				c.cmdb.RetCh <- &sucmd.CommandRet{nil, err}
			}
		}
	}
}

// Cmd takes in a command struct generated from a CommandFactory and executes
// it on the StorageUnit named by unitname
func Cmd(unitname *string, cmd sucmd.Command) (interface{},error) {
	cmdb := sucmd.CommandBundle{
		Cmd: cmd,
		RetCh: make(chan *sucmd.CommandRet),
	}
	cmdw := commandWrap{
		unitname: unitname,
		cmdb: &cmdb,
	}
	select {
	case sm.cmdCh <- &cmdw:
	case <-time.After(10 * time.Second):
		return nil, errors.New("Contacting storage manager timedout")
	}

	select {
	case ret := <- cmdb.RetCh:
		return ret.Ret, ret.Err
	case <-time.After(10 * time.Second):
		return nil, errors.New("Timeout receiving command results")
	}
}
