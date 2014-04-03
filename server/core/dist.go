package core

import (
	"github.com/grooveshark/golib/gslog"

	"github.com/mediocregopher/hyrax/server/config"
	"github.com/mediocregopher/hyrax/server/core/keychanges"
	"github.com/mediocregopher/hyrax/server/dist"
	"github.com/mediocregopher/hyrax/types"
)

// Manager for connections to other nodes we are pulling global events from
var PullFromGlobalManager = dist.New("MGLOBAL")

// Manager for connections to other nodes we are pulling local events from
// (these come from other nodes calling ALISTENTOME)
var PullFromLocalManager = dist.New("MLOCAL")

// Manager for connection to other nodes we are calling ALISTENTOME on,
// effectively commanding them to pull local events from us
var PushToManager = dist.New("ALISTENTOME", config.MyEndpoint.String())

func init() {
	for i := 0; i < 20; i++ {
		go clusterSpin()
	}
}

func clusterSpin() {
	var err error
	var a *types.Action
	for {
		err = nil
		a = nil
		select {
		case a = <-PullFromGlobalManager.PushCh:
			gslog.Debugf("Got %v from global", a)
			err = keychanges.PubGlobal(a)
		case a = <-PullFromLocalManager.PushCh:
			gslog.Debugf("Got %v from local", a)
			err = keychanges.PubGlobal(a)
		case _ = <-PushToManager.PushCh:
		}

		if err != nil {
			gslog.Errorf("PubGlobal(%v) got %s", a, err)
		}
	}
}

// Reads the cluster information from the config and attempts to set it up. If
// this isn't the first time this function has been called it will close all
// existing cluster connections and make new ones
func Clusterize() error {
	// TODO make a bit more resilient to errors, if we encounter one we want to
	// send it back but not disconinue execution

	err := resetManager(
		PullFromGlobalManager,
		config.PullFromEndpoints,
		"MGLOBAL",
	)
	if err != nil {
		return err
	}

	err = resetManager(
		PushToManager,
		config.PushToEndpoints,
		"ALISTENTOME",
		config.MyEndpoint.String(),
	)
	if err != nil {
		return err
	}

	return nil
}

func resetManager(
	m *dist.Manager,
	endpoints []types.ListenEndpoint,
	cmd string, args ...interface{}) error {

	if err := m.CloseAll(); err != nil {
		return err
	}
	m.SetCommand(cmd, args...)
	for _, le := range endpoints {
		if err := m.EnsureClient(le.String()); err != nil {
			return err
		}
	}
	return nil
}
