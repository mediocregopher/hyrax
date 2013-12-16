package main

import (
	"github.com/mediocregopher/hyrax/server/auth"
	"github.com/mediocregopher/hyrax/server/config"
	"github.com/mediocregopher/hyrax/server/dist"
	"github.com/mediocregopher/hyrax/server/net"
	"github.com/mediocregopher/hyrax/server/storage-router"
	"github.com/mediocregopher/hyrax/translate"
	"strings"
	"log"
)

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	if config.FirstNode {
		secrets := config.InitSecrets
		log.Println("This is the first node, loading up the secrets")
		for _, secret := range secrets {
			log.Println("Loading secret:", string(secret))
			auth.AddGlobalSecret(secret)
		}

		storageAddr := config.StorageAddr
		log.Println("Connecting to storage unit at", storageAddr)
		if err := router.SetBucket(0, "tcp", storageAddr); err != nil {
			log.Fatal(err)
		}
	}

	meshListenAddr := config.MeshListenAddr
	meshAdvertiseAddr := config.MeshAdvertiseAddr
	log.Println("Creating mesh listener at", meshListenAddr)
	if err = dist.Init(meshListenAddr); err != nil {
		log.Fatal(err)
	} else if err = dist.AddNode(&meshAdvertiseAddr); err != nil {
		log.Fatal(err)
	}

	listens := config.ListenAddrs
	for i := range listens {
		log.Println("Listening for clients at", listens[i])
		go listenHandler(&listens[i])
	}

	select {}
}

func listenHandler(l *config.ListenAddr) {
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
