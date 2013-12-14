package storage

import (
	"time"
	"errors"
	"github.com/mediocregopher/hyrax/server/storage-router/storage/redis"
	"github.com/mediocregopher/hyrax/server/storage-router/storage/unit"
	"github.com/mediocregopher/hyrax/server/storage-router/storage/command"
)

const UNITSIZE = 10

type LocatorFunc func([]*unit.StorageUnit) *unit.StorageUnit

type commandWrap struct {
	unitname *string
	cmdb *command.CommandBundle
}

type addUnitWrap struct {
	uname string
	u *unit.StorageUnit
}

type storageManager struct {
	units map[string]*unit.StorageUnit
	addUnitCh chan *addUnitWrap
	remUnitCh chan string
	cmdCh chan *commandWrap
}

var sm = storageManager{
	units: map[string]*unit.StorageUnit{},
	addUnitCh: make(chan *addUnitWrap),
	remUnitCh: make(chan string),
	cmdCh: make(chan *commandWrap),
}

var NewStorageUnitConn = redis.New

func init() {
	go sm.spin()
}

func (sm *storageManager) spin() {
	for {
		select {
		case uwrap := <-sm.addUnitCh:
			if _, ok := sm.units[uwrap.uname]; !ok {
				sm.units[uwrap.uname] = uwrap.u
			}

		case uname := <-sm.remUnitCh:
			if u, ok := sm.units[uname]; ok {
				u.Close()
				delete(sm.units, uname)
			}

		case c := <-sm.cmdCh:
			if unit, ok := sm.units[*c.unitname]; ok {
				//We assume our StorageUnits won't block on this call. If any do
				//we have big problems anyway
				unit.Cmd(c.cmdb)
			} else {
				err := errors.New("Unknown storage unit: "+*c.unitname)
				c.cmdb.RetCh <- &command.CommandRet{nil, err}
			}
		}
	}
}

// AddUnit adds a StorageUnit with UNITSIZE StorageUnitConns to the manager
// under the given name
func AddUnit(name, conntype, addr string, extra ...interface{}) error {
	conns := make([]unit.StorageUnitConn, UNITSIZE)
	for i := range conns {
		conns[i] = NewStorageUnitConn()
	}
	
	unit, err := unit.NewStorageUnit(conns, conntype, addr, extra...)
	if err != nil {
		return err
	}

	add := addUnitWrap{name, unit}
	sm.addUnitCh <- &add
	return nil
}

// RemUnit closes and removes a StorageUnit of the given name from the manager
// (assuming it existed in the first place)
func RemUnit(name string) {
	sm.remUnitCh <- name
}

// Cmd takes in a command struct generated from a CommandFactory and executes
// it on the StorageUnit named by unitname
func Cmd(unitname *string, cmd command.Command) (interface{},error) {
	cmdb := command.CommandBundle{
		Cmd: cmd,
		RetCh: make(chan *command.CommandRet),
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
