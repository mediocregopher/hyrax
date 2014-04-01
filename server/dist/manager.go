package dist

import (
	"github.com/grooveshark/golib/gslog"
	"time"

	"github.com/mediocregopher/hyrax/client"
	"github.com/mediocregopher/hyrax/types"
)

type call struct {
	listenEndpoint string
	retCh          chan error
}

type setCmdCall struct {
	cmd  string
	args []interface{}
}

// Manages connections to other hyrax nodes which will perform some command
// periodically. If the connections are cut it will attempt to reconnect them
// periodically as well.
type Manager struct {
	clients map[string]*managerClient

	// All push messages on any clients being managed will be pused down this
	// channel
	PushCh chan *types.ClientCommand
	cmd    string
	args   []interface{}
	period time.Duration

	ensureCh   chan *call
	setCmdCh   chan *setCmdCall
	closeCh    chan *call
	closeAllCh chan *call
}

type managerClient struct {
	le      *types.ListenEndpoint
	cl      client.Client
	pushCh  chan *types.ClientCommand
	closeCh chan struct{}
}

func New(cmd string, args ...interface{}) *Manager {
	m := Manager{
		clients:    map[string]*managerClient{},
		PushCh:     make(chan *types.ClientCommand),
		cmd:        cmd,
		args:       args,
		period:     5 * time.Second,
		ensureCh:   make(chan *call),
		setCmdCh:   make(chan *setCmdCall),
		closeCh:    make(chan *call),
		closeAllCh: make(chan *call),
	}
	go m.spin()
	return &m
}

// Takes in the listen address, which is the same as that given in
// server/config. Ensures there is a client connected to that address which is
// periodically calling the manager's command
func (m *Manager) EnsureClient(listenEndpoint string) error {
	c := call{listenEndpoint, make(chan error)}
	m.ensureCh <- &c
	return <-c.retCh
}

// Tells the manager to change what command it is periodically sending to the
// other nodes in the cluster
func (m *Manager) SetCommand(cmd string, args ...interface{}) {
	m.setCmdCh <- &setCmdCall{cmd, args}
}

// Closes the connection to the given listenEndpoint (see EnsureClient)
func (m *Manager) CloseClient(listenEndpoint string) error {
	c := call{listenEndpoint, make(chan error)}
	m.closeCh <- &c
	return <-c.retCh
}

// Closes all connections
func (m *Manager) CloseAll() error {
	c := call{"", make(chan error)}
	m.closeAllCh <- &c
	return <-c.retCh
}

func (m *Manager) spin() {
	for {
		select {
		case c := <-m.ensureCh:
			c.retCh <- m.ensureClient(c.listenEndpoint)
		case c := <-m.setCmdCh:
			m.cmd = c.cmd
			m.args = c.args
		case c := <-m.closeCh:
			c.retCh <- m.closeClient(c.listenEndpoint)
		case c := <-m.closeAllCh:
			c.retCh <- m.closeAll()
		}
	}
}

func (m *Manager) ensureClient(listenEndpoint string) error {
	le, err := types.ListenEndpointFromString(listenEndpoint)
	if err != nil {
		return err
	}

	if _, ok := m.clients[listenEndpoint]; ok {
		return nil
	}

	pushCh := make(chan *types.ClientCommand)
	cl, err := client.NewClient(le, pushCh)
	if err != nil {
		return err
	}

	mcl := managerClient{
		le:      le,
		cl:      cl,
		pushCh:  pushCh,
		closeCh: make(chan struct{}),
	}
	m.clients[listenEndpoint] = &mcl
	go m.clientSpin(&mcl)
	return nil
}

func (m *Manager) closeClient(listenEndpoint string) error {
	if mcl, ok := m.clients[listenEndpoint]; ok {
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
			cmd := client.CreateClientCommand(m.cmd, "", "", "", m.args...)
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
			cmd := client.CreateClientCommand(m.cmd, "", "", "")
			if _, err := mcl.cl.Cmd(cmd); err != nil {
				gslog.Errorf("managerClient Cmd(%v): %s", cmd, err)
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
			cl, err := client.NewClient(mcl.le, mcl.pushCh)
			if err != nil {
				gslog.Errorf("Error reconnecting to %s: %s", mcl.le, err)
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
