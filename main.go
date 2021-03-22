package main

import (
	"github.com/joho/godotenv"
	"github.com/sharpvik/log-go/v2"

	"github.com/sharpvik/ava/configs"
	"github.com/sharpvik/ava/server"
)

func init() {
	log.SetLevel(log.LevelDebug)
	if err := godotenv.Load(); err != nil {
		log.Error(err)
	}
}

func main() {
	config := configs.MustInit()
	log.Debug("init successfull")

	serv := server.NewServer(config.Server)
	done := make(chan bool, 1)
	go serv.ServeWithGrace(done)

	<-done
	log.Debug("server stopped")
}
