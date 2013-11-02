package main

import (
	"github.com/mediocregopher/hyrax/server/auth"
	"github.com/mediocregopher/hyrax/server/config"
	"github.com/mediocregopher/hyrax/server/dist"
	"github.com/mediocregopher/hyrax/server/net"
	"github.com/mediocregopher/hyrax/server/storage-router/storage"
	"github.com/mediocregopher/hyrax/translate"
)

func main() {
	err := config.Load()
	if err != nil {
		panic(err)
	}

	if config.FirstNode {
		secrets := config.InitSecrets
		for _, secret := range secrets {
			auth.AddGlobalSecret(secret)
		}
	}

	storageAddr := config.StorageAddr
	err = storage.AddUnit(storageAddr, "tcp", storageAddr)
	if err != nil {
		panic(err)
	}

	meshListenAddr := config.MeshListenAddr
	meshAdvertiseAddr := config.MeshAdvertiseAddr
	if err = dist.Init(meshListenAddr); err != nil {
		panic(err)
	} else if err = dist.AddNode(&meshAdvertiseAddr); err != nil {
		panic(err)
	}

	listens := config.ListenAddrs
	for i := range listens {
		go listenHandler(&listens[i])
	}

	select {}
}

func listenHandler(l *config.ListenAddr) {
	var trans translate.Translator
	switch l.Format {
	case config.LFORMAT_JSON:
		trans = &translate.JsonTranslator{}
	}

	switch l.Type {
	case  config.LTYPE_TCP:
		if err := net.TcpListen(l.Addr, trans); err != nil {
			panic(err)
		}
	}
}
