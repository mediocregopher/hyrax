package dist

import (
	"github.com/grooveshark/golib/gslog"
	"time"

	"github.com/mediocregopher/hyrax/client"
	"github.com/mediocregopher/hyrax/server/config"
	"github.com/mediocregopher/hyrax/types"
)

type call struct {
	listenEndpoint *types.ListenEndpoint
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
	PushCh  chan *types.Action
	cmd     string
	args    []interface{}
	period  time.Duration
	timeout time.Duration

	ensureCh   chan *call
	setCmdCh   chan *setCmdCall
	closeCh    chan *call
	closeAllCh chan *call
	getAllCh   chan chan []*types.ListenEndpoint
}

type managerClient struct {
	le      *types.ListenEndpoint
	cl      client.Client
	pushCh  chan *types.Action
	closeCh chan struct{}
	touchCh chan struct{}
}

func newMan(cmd string, args ...interface{}) *Manager {
	return &Manager{
		clients:    map[string]*managerClient{},
		PushCh:     make(chan *types.Action),
		cmd:        cmd,
		args:       args,
		period:     5 * time.Second,
		ensureCh:   make(chan *call),
		setCmdCh:   make(chan *setCmdCall),
		closeCh:    make(chan *call),
		closeAllCh: make(chan *call),
		getAllCh:   make(chan chan []*types.ListenEndpoint),
	}
}

func New(cmd string, args ...interface{}) *Manager {
	m := newMan(cmd, args...)
	go m.spin()
	return m
}

// Like New, but the manager will timeout client connections if they aren't
// ensured every t units of time
func NewTimeout(t time.Duration, cmd string, args ...interface{}) *Manager {
	m := newMan(cmd, args...)
	m.timeout = t
	go m.spin()
	return m
}

// Takes in the listen address, which is the same as that given in
// server/config. Ensures there is a client connected to that address which is
// periodically calling the manager's command
func (m *Manager) EnsureClient(le *types.ListenEndpoint) error {
	c := call{le, make(chan error)}
	m.ensureCh <- &c
	return <-c.retCh
}

// Tells the manager to change what command it is periodically sending to the
// other nodes in the cluster
func (m *Manager) SetCommand(cmd string, args ...interface{}) {
	m.setCmdCh <- &setCmdCall{cmd, args}
}

// Closes the connection to the given listenEndpoint (see EnsureClient)
func (m *Manager) CloseClient(le *types.ListenEndpoint) error {
	c := call{le, make(chan error)}
	m.closeCh <- &c
	return <-c.retCh
}

// Closes all connections
func (m *Manager) CloseAll() error {
	c := call{nil, make(chan error)}
	m.closeAllCh <- &c
	return <-c.retCh
}

// Returns all currently active endpoints
func (m *Manager) GetAll() []*types.ListenEndpoint {
	retCh := make(chan []*types.ListenEndpoint)
	m.getAllCh <- retCh
	return <-retCh
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
		case retCh := <-m.getAllCh:
			retCh <- m.getAllClients()
		}
	}
}

func (m *Manager) ensureClient(le *types.ListenEndpoint) error {
	leStr := le.String()
	gslog.Debugf("Ensuring %s connection to node %s", m.cmd, leStr)

	if mcl, ok := m.clients[leStr]; ok {
		mcl.touchCh <- struct{}{}
		return nil
	}

	pushCh := make(chan *types.Action)
	cl, err := client.NewClient(le, pushCh)
	if err != nil {
		return err
	}
	gslog.Debugf("Created %s connection to %s", m.cmd, leStr)

	mcl := managerClient{
		le:      le,
		cl:      cl,
		pushCh:  pushCh,
		closeCh: make(chan struct{}),
		touchCh: make(chan struct{}),
	}
	m.clients[leStr] = &mcl
	go m.clientSpin(&mcl)
	return nil
}

func (m *Manager) closeClient(le *types.ListenEndpoint) error {
	leStr := le.String()
	gslog.Debugf("Closing %s connection to node %s", m.cmd, leStr)
	if mcl, ok := m.clients[leStr]; ok {
		close(mcl.closeCh)
		delete(m.clients, leStr)
	}
	return nil
}

func (m *Manager) closeAll() error {
	for leStr, mcl := range m.clients {
		close(mcl.closeCh)
		delete(m.clients, leStr)
	}
	return nil
}

func (m *Manager) getAllClients() []*types.ListenEndpoint {
	ret := make([]*types.ListenEndpoint, 0, len(m.clients))
	for leStr := range m.clients {
		le, err := types.ListenEndpointFromString(leStr)
		if err != nil {
			gslog.Errorf("dist.Manager ListenEndpointFromString: %s", err)
			continue
		}
		ret = append(ret, le)
	}
	return ret
}

func (m *Manager) clientSpin(mcl *managerClient) {
	var timeout *time.Timer
	var timeoutCh <-chan time.Time
	if m.timeout > 0 {
		timeout = time.NewTimer(m.timeout)
		timeoutCh = timeout.C
	}

	ticker := time.NewTicker(m.period)
	doCmd := true

spinloop:
	for {

		if doCmd {
			secret := config.InteractionSecret
			cmd := client.CreateAction(m.cmd, "", "", secret, m.args...)
			if _, err := mcl.cl.Cmd(cmd); err != nil {
				gslog.Errorf("dist cmd %s: %s", m.cmd, err)
				mcl.cl.Close()
				if !mcl.resurrect() {
					break spinloop
				} else {
					continue
				}
			} else {
				doCmd = false
			}
		}

		select {
		case a, ok := <-mcl.pushCh:
			if !ok {
				break spinloop
			}
			m.PushCh <- a
		case <-timeoutCh:
			gslog.Warnf("Closing %s connection to %s", m.cmd, mcl.le)
			if err := m.CloseClient(mcl.le); err != nil {
				gslog.Errorf("Closing %s conn to %s: %s", m.cmd, mcl.le, err)
			}
			break spinloop
		case <-mcl.touchCh:
			timeout.Reset(m.timeout)
		case <-mcl.closeCh:
			break spinloop
		case <-ticker.C:
			doCmd = true
		}
	}

	mcl.cl.Close()
	ticker.Stop()
	if m.period > 0 {
		timeout.Stop()
	}
}

func (mcl *managerClient) resurrect() bool {
	clCh := make(chan client.Client)

	go func() {
		for {
			time.Sleep(2 * time.Second)
			gslog.Debug("Going to resurrect new client on %s", mcl.le)
			cl, err := client.NewClient(mcl.le, mcl.pushCh)
			if err != nil {
				gslog.Errorf("Error reconnecting to %s: %s", mcl.le, err)
				continue
			}
			clCh <- cl
			return
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
