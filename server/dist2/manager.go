package dist2

import (
	"sync"
	
	"github.com/mediocregopher/hyrax/client"
	"github.com/mediocregopher/hyrax/server/config"
	"github.com/mediocregopher/hyrax/types"
)

var monClients = map[string]client.Client{}
var monClientsLock sync.RWMutex

// When pushed key changes come in from other nodes they get put on here
var ClientCommandCh = make(chan *types.ClientCommand)

// Makes sure there is a connection monitoring all keychanges on the given
// server. Takes in the listen address, which is the same as that given in
// server/config.
func EnsureMonClient(listenAddr string) error {
	la, err := config.ParseListenAddr(listenAddr)
	if err != nil {
		return err
	}

	monClientsLock.Lock()
	defer monClientsLock.Unlock()
	if _, ok := monClients[listenAddr]; ok {
		return nil
	}

	pushCh := make(chan *types.ClientCommand)
	cl, err := client.NewClient(la.Format, la.Type, la.Addr, pushCh)
	if err != nil {
		return err
	}

	// TODO secret
	cmd := client.CreateClientCommand([]byte("MALL"), nil, nil, nil)
	if _, err = cl.Cmd(cmd); err != nil {
		cl.Close()
		return err
	}

	monClients[listenAddr] = cl
	return nil
}

// Stops receiving keychanges from the given listening address, if it was
// previously
func StopClient(listenAddr string) error {
	monClientsLock.Lock()
	defer monClientsLock.Unlock()
	if client, ok := monClients[listenAddr]; ok {
		client.Close()
		delete(monClients, listenAddr)
	}
	return nil
}

func pushSpin(pushCh chan *types.ClientCommand) {
	for cmd := range pushCh {
		ClientCommandCh <- cmd
	}
}
