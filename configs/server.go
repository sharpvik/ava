package configs

import (
	"flag"
	"net/http"
	"os"

	"github.com/sharpvik/log-go/v2"
)

type Server struct {
	StorageDir http.Dir
	DevMode    bool
	APIKey     string
}

func mustInitServer() *Server {
	log.Debug("config server")
	storageDir := flag.String("dir", "storage", "specify custom storage folder")
	devMode := flag.Bool("dev", false, "run server in development mode")
	apiKey := os.Getenv("API_KEY")
	flag.Parse()
	return &Server{
		StorageDir: http.Dir(*storageDir),
		DevMode:    *devMode,
		APIKey:     apiKey,
	}
}
