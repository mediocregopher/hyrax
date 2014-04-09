package dist

import (
	"github.com/grooveshark/golib/gslog"
	"time"

	"github.com/mediocregopher/hyrax/server/config"
	"github.com/mediocregopher/hyrax/server/core/keychanges"
	"github.com/mediocregopher/hyrax/server/dist"
	"github.com/mediocregopher/hyrax/types"
)

// Manager for connections to other nodes we are pulling global events from
var PullFromGlobalManager = dist.New("MGLOBAL")

// Manager for connections to other nodes we are pulling local events from
// (these come from other nodes calling ALISTENTOME)
var PullFromLocalManager = dist.NewTimeout(10*time.Second, "MLOCAL")

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
// this isn't the first time this function has been called it will do a diff and
// open/close whatever connections are needed, and leave the remaining ones
// untouched.
func Clusterize() error {
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
	les []*types.ListenEndpoint,
	cmd string, args ...interface{}) error {

	m.SetCommand(cmd, args...)
	oldLes := m.GetAll()

	lesM := map[string]bool{}
	oldLesM := map[string]bool{}
	for i := range les {
		lesM[les[i].String()] = true
	}
	for i := range oldLes {
		oldLesM[oldLes[i].String()] = true
	}

	// Add loop
	for _, le := range les {
		if _, ok := oldLesM[le.String()]; ok {
			continue
		}
		gslog.Infof("Adding %s connection to %s", cmd, le)
		if err := m.EnsureClient(le); err != nil {
			return err
		}
	}

	// Remove loop
	for _, oldLe := range oldLes {
		if _, ok := lesM[oldLe.String()]; ok {
			continue
		}
		gslog.Infof("Removing %s connection to %s", cmd, oldLe)
		if err := m.CloseClient(oldLe); err != nil {
			return err
		}
	}

	return nil
}
