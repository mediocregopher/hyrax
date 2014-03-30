package dist2

import (
	"github.com/mediocregopher/hyrax/client"
	"github.com/mediocregopher/hyrax/types"
	"time"
)

type call struct {
	listenAddr string
	retCh      chan error
}

// Manages connections to other hyrax nodes which will perform some command
// periodically. If the connections are cut it will attempt to reconnect them
// periodically as well.
type Manager struct {
	clients map[string]*managerClient

	// All push messages on any clients being managed will be pused down this
	// channel
	PushCh  chan *types.ClientCommand
	cmd     []byte
	period  time.Duration

	ensureCh   chan *call
	closeCh    chan *call
	closeAllCh chan *call
}

type managerClient struct {
	la       *types.ListenAddr
	cl       client.Client
	pushCh   chan *types.ClientCommand
	closeCh  chan struct{}
}

func New(cmd string) *Manager {
	m := Manager{
		clients:    map[string]*managerClient{},
		PushCh:     make(chan *types.ClientCommand),
		cmd:        []byte(cmd),
		period:     5 * time.Second,
		ensureCh:   make(chan *call),
		closeCh:    make(chan *call),
		closeAllCh: make(chan *call),
	}
	go m.spin()
	return &m
}

// Takes in the listen address, which is the same as that given in
// server/config. Ensures there is a client connected to that address which is
// periodically calling the manager's command
func (m *Manager) EnsureClient(listenAddr string) error {
	c := call{listenAddr, make(chan error)}
	m.ensureCh <- &c
	return <-c.retCh
}

// Closes the connection to the given listenAddr (see EnsureClient)
func (m *Manager) CloseClient(listenAddr string) error {
	c := call{listenAddr, make(chan error)}
	m.closeCh <- &c
	return <-c.retCh
}

// Closes all connections
func (m *Manager) CloseAll() error {
	c := call{"", make(chan error)}
	m.closeAllCh <- &c
	return <-c.retCh
}

func (m* Manager) spin() {
	for {
		select {
		case c := <-m.ensureCh:
			c.retCh <- m.ensureClient(c.listenAddr)
		case c := <-m.closeCh:
			c.retCh <- m.closeClient(c.listenAddr)
		case c := <-m.closeAllCh:
			c.retCh <-m.closeAll()
		}
	}
}

func (m *Manager) ensureClient(listenAddr string) error {
	la, err := types.ListenAddrFromString(listenAddr)
	if err != nil {
		return err
	}

	if _, ok := m.clients[listenAddr]; ok {
		return nil
	}

	pushCh := make(chan *types.ClientCommand)
	cl, err := client.NewClient(la, pushCh)
	if err != nil {
		return err
	}

	mcl := managerClient{
		la:      la,
		cl:      cl,
		pushCh:  pushCh,
		closeCh: make(chan struct{}),
	}
	m.clients[listenAddr] = &mcl
	go m.clientSpin(&mcl)
	return nil
}

func (m *Manager) closeClient(listenAddr string) error {
	if mcl, ok := m.clients[listenAddr]; ok {
		close(mcl.closeCh)
	}
	return nil
}

func (m *Manager) closeAll() error {
	for _, mcl := range m.clients {
		close(mcl.closeCh)
	}
	return nil
}


func (m *Manager) clientSpin(mcl *managerClient) {
	ticker := time.NewTicker(m.period)

spinloop:
	for {
		select {
		case cc, ok := <-mcl.pushCh:
			if !ok {
				break spinloop
			}
			m.PushCh <- cc
		case <-ticker.C:
			// TODO secret
			cmd := client.CreateClientCommand(m.cmd, nil, nil, nil)
			if _, err := mcl.cl.Cmd(cmd); err != nil {
				mcl.cl.Close()
				break spinloop
			}
		}
	}

	ticker.Stop()
}

func (mcl *managerClient) spin(m *Manager) {
	ticker := time.NewTicker(m.period)
	doCmd := true
	resurrect := false

spinloop:
	for {
		if doCmd {
			// TODO secret
			cmd := client.CreateClientCommand(m.cmd, nil, nil, nil)
			if _, err := mcl.cl.Cmd(cmd); err != nil {
				// TODO log error
				resurrect = true
			}
			doCmd = false
		}

		if !resurrect {
			select {
			case <-mcl.closeCh:
				break spinloop
			case cc, ok := <-mcl.pushCh:
				if !ok {
					resurrect = true
					break
				}
				m.PushCh <- cc
			case <-ticker.C:
				doCmd = true
			}
		}

		if resurrect {
			mcl.cl.Close()
			if !mcl.resurrect() {
				break spinloop
			}
		}
	}

	ticker.Stop()
	mcl.cl.Close()
}

func (mcl *managerClient) resurrect() bool {
	clCh := make(chan client.Client)

	go func() {
		for {
			cl, err := client.NewClient(mcl.la, mcl.pushCh)
			if err != nil {
				// TODO log error
				time.Sleep(2 * time.Second)
				continue
			}
			clCh <- cl
		}
	}()

	for {
		select {
		case cl := <-clCh:
			mcl.cl = cl
			return true
		case <-mcl.closeCh:
			return false
		}
	}
}
