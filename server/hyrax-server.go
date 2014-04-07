package main

import (
	"github.com/grooveshark/golib/gslog"
	"os"
	"os/signal"
	"syscall"

	"github.com/mediocregopher/hyrax/server/core"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR1)
	for _ = range c {
		if err := core.Reload(); err != nil {
			gslog.Errorf("reload: %s", err)
		}
	}
}
