package main

import (
	"github.com/grooveshark/golib/gslog"
	"strings"

	"github.com/mediocregopher/hyrax/server/auth"
	"github.com/mediocregopher/hyrax/server/config"
	"github.com/mediocregopher/hyrax/server/core"
	"github.com/mediocregopher/hyrax/server/net"
	"github.com/mediocregopher/hyrax/translate"
	"github.com/mediocregopher/hyrax/types"
)

func main() {
	secrets := config.InitSecrets
	gslog.Info("Loading up the secrets")
	for _, secret := range secrets {
		gslog.Info("Loading secret:", secret)
		auth.AddGlobalSecret(secret)
	}

	gslog.Infof("Connecting to datastore at: %s", config.StorageInfo)
	if err := core.SetupStorage(); err != nil {
		gslog.Fatal(err.Error())
	}

	listens := config.ListenEndpoints
	for i := range listens {
		gslog.Infof("Listening for clients at: %s", listens[i])
		go listenHandler(&listens[i])
	}

	gslog.Info("Connecting to other nodes in the cluster")
	if err := core.Clusterize(); err != nil {
		gslog.Fatal(err.Error())
	}

	select {}
}

func listenHandler(l *types.ListenEndpoint) {
	trans, err := translate.StringToTranslator(l.Format)
	if err != nil {
		gslog.Fatal(err.Error())
	}

	switch strings.ToLower(l.Type) {
	case "tcp":
		if err := net.TcpListen(l.Addr, trans); err != nil {
			gslog.Fatal(err.Error())
		}
	}
}
