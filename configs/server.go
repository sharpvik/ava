package configs

import (
	"flag"
	"net/http"

	"github.com/sharpvik/log-go/v2"
)

type Server struct {
	StorageDir http.Dir
	DevMode    bool
}

func mustInitServer() Server {
	log.Debug("config server")
	dev, dir := parseFlags()
	return Server{
		StorageDir: http.Dir(dir),
		DevMode:    dev,
	}
}

func parseFlags() (dev bool, dir string) {
	devMode := flag.Bool("dev", false, "run server in development mode")
	storageDir := flag.String("dir", "storage", "specify custom storage folder")
	flag.Parse()
	return *devMode, *storageDir
}
