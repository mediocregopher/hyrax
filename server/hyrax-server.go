package main

import (
	"github.com/grooveshark/golib/gslog"

	"github.com/mediocregopher/hyrax/server/core"
)

func main() {

	if err := core.Configure(); err != nil {
		gslog.Fatal(err.Error())
	}

	select {}
}

