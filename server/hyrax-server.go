package main

import (
	"github.com/mediocregopher/hyrax/server/auth"
	"github.com/mediocregopher/hyrax/server/config"
	"github.com/mediocregopher/hyrax/server/core"
	"github.com/mediocregopher/hyrax/server/net"
	"github.com/mediocregopher/hyrax/server/storage-router"
	"github.com/mediocregopher/hyrax/translate"
	"github.com/mediocregopher/hyrax/types"
	"log"
	"strings"
)

func main() {
	secrets := config.InitSecrets
	log.Println("Loading up the secrets")
	for _, secret := range secrets {
		log.Println("Loading secret:", string(secret))
		auth.AddGlobalSecret(secret)
	}

	storageAddr := config.StorageAddr
	log.Println("Connecting to storage unit at", storageAddr)
	if err := router.SetBucket(0, "tcp", storageAddr); err != nil {
		log.Fatal(err)
	}

	listens := config.ListenEndpoints
	for i := range listens {
		log.Println("Listening for clients at", listens[i])
		go listenHandler(&listens[i])
	}

	log.Println("Connecting to other nodes in the cluster")
	if err := core.Clusterize(); err != nil {
		log.Fatal(err)
	}

	select {}
}

func listenHandler(l *types.ListenEndpoint) {
	trans, err := translate.StringToTranslator(l.Format)
	if err != nil {
		log.Fatal(err)
	}

	switch strings.ToLower(l.Type) {
	case "tcp":
		if err := net.TcpListen(l.Addr, trans); err != nil {
			log.Fatal(err)
		}
	}
}
