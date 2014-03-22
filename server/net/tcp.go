package net

import (
	"bufio"
	"log"
	"github.com/mediocregopher/manatcp"
	"time"

	"github.com/mediocregopher/hyrax/server/client"
	crouter "github.com/mediocregopher/hyrax/server/client-router"
	stypes "github.com/mediocregopher/hyrax/server/types"
	"github.com/mediocregopher/hyrax/types"
	"github.com/mediocregopher/hyrax/translate"
)

type tcpListener struct {
	trans translate.Translator
}

func (tl *tcpListener) Connected (
	lc *manatcp.ListenerConn) (manatcp.ServerClient, bool) {

	cid := client.NewClientId()
	c := tcpClient{
		cmdPushCh: make(chan *types.ClientCommand),
		lconn:     lc,
		id:        cid,
		trans:     tl.trans,
	}
	if err := crouter.Add(&c); err != nil {
		log.Printf("tcpListener got %s adding to crouter", err) 
		return nil, true
	}

	go c.pushProxy()
	return &c, false	
}

func TcpListen(addr string, trans translate.Translator) error {
	_, err := manatcp.Listen(&tcpListener{trans}, addr)
	return err
}

type tcpClient struct {
	cmdPushCh chan *types.ClientCommand
	lconn     *manatcp.ListenerConn
	id        stypes.ClientId
	trans     translate.Translator
}

func (tc *tcpClient) pushProxy() {
	for cmd := range tc.cmdPushCh {
		tc.lconn.PushCh <- cmd
	}
}

func (tc *tcpClient) ClientId() stypes.ClientId {
	return tc.id
}

func (tc *tcpClient) CommandPushChannel() chan<- *types.ClientCommand {
	return tc.cmdPushCh
}

func (tc *tcpClient) Read(buf *bufio.Reader) (interface{}, bool) {
	b, err := buf.ReadBytes('\n')
	return b, err != nil
}

func (tc *tcpClient) Write(buf *bufio.Writer, clientRet interface{}) bool {
	b, err := tc.trans.FromClientReturn(clientRet.(*types.ClientReturn))
	if err != nil {
		log.Printf("tcpClient got %s from FromClientReturn(%v)", err, clientRet)
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
	cc, err := tc.trans.ToClientCommand(cmdRaw.([]byte))
	if err != nil {
		return types.ErrorReturn(err), true, false
	}
	return client.RunCommand(tc.id, cc), true, false
}

func (tc *tcpClient) Closing() {
	client.ClientClosed(tc.id)
	crouter.RemByClient(tc)
	// We sleep some seconds just in case anything is still pushing to the
	// command channel
	time.Sleep(5 * time.Second)
	close(tc.cmdPushCh)
}
