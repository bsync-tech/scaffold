package main

import (
	"github.com/bsync-tech/mlog"
	"github.com/bsync-tech/scaffold/config"
	"github.com/bsync-tech/scaffold/http/server"
	"github.com/bsync-tech/scaffold/log"
)

func main() {
	defer mlog.Sync()

	log.New(log.WithPath("conf/log.yaml"))
	log.Initialize()

	config.New(config.WithPath("conf/config.yaml"))
	config.Initialize()

	server.HttpServerRun()

	// quit := make(chan os.Signal)
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// <-quit
}
