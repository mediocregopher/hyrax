package core

import (
	"github.com/grooveshark/golib/gslog"
	"strings"

	"github.com/mediocregopher/hyrax/server/auth"
	"github.com/mediocregopher/hyrax/server/config"
	"github.com/mediocregopher/hyrax/server/core/dist"
	"github.com/mediocregopher/hyrax/server/listen"
	"github.com/mediocregopher/hyrax/translate"
	"github.com/mediocregopher/hyrax/types"
)

func init() {
	if err := InitialConfigure(); err != nil {
		gslog.Fatal(err.Error())
	}
}

// Does all configuration for the hyrax node
func InitialConfigure() error {
	if err := gslog.SetMinimumLevel(config.LogLevel); err != nil {
		return err
	}
	gslog.Infof("Setting logging point to %s", config.LogFile)
	if err := gslog.SetLogFile(config.LogFile); err != nil {
		return err
	}

	secrets := config.InitSecrets
	gslog.Info("Loading up the secrets")
	for _, secret := range secrets {
		auth.AddGlobalSecret(secret)
	}

	if err := SetupStorage(); err != nil {
		return err
	}

	listens := config.ListenEndpoints
	for i := range listens {
		if err := listenHandler(listens[i]); err != nil {
			return err
		}
	}

	if err := dist.Clusterize(); err != nil {
		return err
	}

	return nil
}

// Performs whatever actions are supported on reload
func Reload() error {
	if err := config.Load(); err != nil {
		return err
	}
	return dist.Clusterize()
}

func listenHandler(l *types.ListenEndpoint) error {
	trans, err := translate.StringToTranslator(l.Format)
	if err != nil {
		return err
	}

	switch strings.ToLower(l.Type) {
	case "tcp":
		if err := listen.TcpListen(l.Addr, trans); err != nil {
			return err
		}
	}

	return nil
}
