package listen

import (
	"bufio"
	"errors"
	"github.com/grooveshark/golib/gslog"
	"github.com/mediocregopher/manatcp"
	"time"

	stypes "github.com/mediocregopher/hyrax/server/types"
	"github.com/mediocregopher/hyrax/translate"
	"github.com/mediocregopher/hyrax/types"
)

// ActionWrap bundles an action with a channel which can be read from to receive
// the return from that action
type ActionWrap struct {
	Action         *types.Action
	Client         stypes.Client
	ActionReturnCh chan *types.ActionReturn
}

// Whenever a client performs an action it will be wrapped and put on this
// channel, waiting for a reply on the ActionReturnCh in the ActionWrap
var ActionWrapCh = make(chan *ActionWrap)

// ClientClosedWrap bundles a client closing event with a channel which will be
// closed when cleanup is completed
type ClientClosedWrap struct {
	Client stypes.Client
	Ch     chan struct{}
}

// Whenever a client is closed it will be put on this channel so that some other
// process can clean it up
var ClientClosedCh = make(chan *ClientClosedWrap)

type tcpListener struct {
	trans translate.Translator
}

func (tl *tcpListener) Connected(
	lc *manatcp.ListenerConn) (manatcp.ServerClient, bool) {

	cid := stypes.NewClientId()
	c := tcpClient{
		cmdPushCh: make(chan *types.Action),
		lconn:     lc,
		id:        cid,
		trans:     tl.trans,
		closeCh:   make(chan struct{}),
	}

	go c.pushProxy()
	return &c, false
}

func TcpListen(addr string, trans translate.Translator) error {
	gslog.Infof("Listening for clients at %s", addr)
	_, err := manatcp.Listen(&tcpListener{trans}, addr)
	return err
}

type tcpClient struct {
	cmdPushCh chan *types.Action
	lconn     *manatcp.ListenerConn
	id        stypes.ClientId
	trans     translate.Translator
	closeCh   chan struct{}
}

func (tc *tcpClient) pushProxy() {
	for cmd := range tc.cmdPushCh {
		tc.lconn.PushCh <- cmd
	}
}

func (tc *tcpClient) ClientId() stypes.ClientId {
	return tc.id
}

func (tc *tcpClient) PushCh() chan<- *types.Action {
	return tc.cmdPushCh
}

func (tc *tcpClient) ClosingCh() <-chan struct{} {
	return tc.closeCh
}

func (tc *tcpClient) Read(buf *bufio.Reader) (interface{}, bool) {
	b, err := buf.ReadBytes('\n')
	return b, err != nil
}

func (tc *tcpClient) Write(buf *bufio.Writer, i interface{}) bool {
	var b []byte
	var err error
	if ar, ok := i.(*types.ActionReturn); ok {
		b, err = tc.trans.FromActionReturn(ar)
	} else if a, ok := i.(*types.Action); ok {
		b, err = tc.trans.FromAction(a)
	} else {
		err = errors.New("invalid type to write")
	}
	if err != nil {
		gslog.Warnf("tcpClient Write(%v): %s", i, err)
		return false
	}
	if _, err := buf.Write(b); err != nil {
		return true
	}
	if _, err := buf.Write([]byte("\n")); err != nil {
		return true
	}
	return false
}

func (tc *tcpClient) HandleCmd(cmdRaw interface{}) (interface{}, bool, bool) {
	a, err := tc.trans.ToAction(cmdRaw.([]byte))
	if err != nil {
		return types.NewActionReturn(err), true, false
	}
	return DispatchAction(tc, a), true, false
}

func (tc *tcpClient) Closing() {
	DispatchClosed(tc)
	// We sleep some seconds just in case anything is still pushing to the
	// command channel
	time.Sleep(5 * time.Second)
	close(tc.cmdPushCh)
	close(tc.closeCh)
}

func DispatchAction(c stypes.Client, a *types.Action) *types.ActionReturn {
	aw := ActionWrap{a, c, make(chan *types.ActionReturn)}
	select {
	case ActionWrapCh <- &aw:
	case <-time.After(5 * time.Second):
		gslog.Error("Timedout sending Action to ActionWrapCh")
		return types.NewActionReturn(errors.New("timeout"))
	}

	select {
	case ar := <-aw.ActionReturnCh:
		return ar
	case <-time.After(5 * time.Second):
		gslog.Error("Timedout receiving ActionReturn from ActionReturnCh")
		return types.NewActionReturn(errors.New("timeout"))
	}
}

func DispatchClosed(c stypes.Client) {
	cc := ClientClosedWrap{c, make(chan struct{})}
	select {
	case ClientClosedCh <- &cc:
	case <-time.After(5 * time.Second):
		gslog.Error("Timedout sending to ClientClosedCh")
		return
	}

	select {
	case <-cc.Ch:
	case <-time.After(5 * time.Second):
		gslog.Error("Timedout waiting for ClientClosedWrap.Ch")
	}
}
