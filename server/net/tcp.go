package net

import (
	"bufio"
	"github.com/grooveshark/golib/gslog"
	"github.com/mediocregopher/manatcp"
	"time"

	"github.com/mediocregopher/hyrax/server/core"
	stypes "github.com/mediocregopher/hyrax/server/types"
	"github.com/mediocregopher/hyrax/translate"
	"github.com/mediocregopher/hyrax/types"
)

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

func (tc *tcpClient) Write(buf *bufio.Writer, ar interface{}) bool {
	b, err := tc.trans.FromActionReturn(ar.(*types.ActionReturn))
	if err != nil {
		gslog.Warnf("tcpClient FromActionReturn(%v): %s", ar, err)
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
	cmd, err := tc.trans.ToAction(cmdRaw.([]byte))
	if err != nil {
		return types.ErrorReturn(err), true, false
	}
	return core.RunCommand(tc, cmd), true, false
}

func (tc *tcpClient) Closing() {
	core.ClientClosed(tc)
	// We sleep some seconds just in case anything is still pushing to the
	// command channel
	time.Sleep(5 * time.Second)
	close(tc.cmdPushCh)
	close(tc.closeCh)
}
