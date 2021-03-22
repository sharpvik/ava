package main

import (
	"github.com/sharpvik/log-go/v2"

	"github.com/sharpvik/ava/configs"
	"github.com/sharpvik/ava/server"
)

// mustInit is responsible for the primary and the most essential initialization
// code that has to run properly.
func mustInit() (config configs.Config) {
	log.SetLevel(log.LevelDebug)
	config = configs.MustInit()
	return
}

func main() {
	config := mustInit()
	log.Debug("init successfull")

	serv := server.NewServer(config.Server)
	done := make(chan bool, 1)
	go serv.ServeWithGrace(done)

	<-done
	log.Debug("server stopped")
}
