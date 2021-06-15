package configs

import (
	"flag"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/sharpvik/log-go/v2"
)

type Server struct {
	StorageDir http.Dir
	DevMode    bool
	APIKey     string
	Extensions []string
	SizeLimit  int64
}

func mustInitServer() *Server {
	log.Debug("config server")
	storageDir := flag.String("dir", "storage", "specify custom storage folder")
	devMode := flag.Bool("dev", false, "run server in development mode")
	apiKey := os.Getenv("AVA_API_KEY")
	extensions := strings.Split(os.Getenv("AVA_ALLOWED_EXTENSIONS"), ",")
	sizeLimit, err := strconv.ParseInt(getOr("AVA_SIZE_LIMIT", "0"), 10, 64)
	if err != nil {
		panic(err)
	}
	flag.Parse()
	return &Server{
		StorageDir: http.Dir(*storageDir),
		DevMode:    *devMode,
		APIKey:     apiKey,
		Extensions: extensions,
		SizeLimit:  sizeLimit,
	}
}

func getOr(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}
