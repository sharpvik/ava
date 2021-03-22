package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/sharpvik/log-go/v2"
	"github.com/sharpvik/mux"
)

func logRequest(pathPrefix string) mux.View {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Infof("%s %s%s", r.Method, pathPrefix, r.URL.String())
	}
}

func newImageName(ext string) string {
	return fmt.Sprintf("%s.%s",
		strconv.FormatInt(time.Now().UnixNano(), 16), ext)
}

func (h *handler) pathToFile(name string) string {
	return path.Join(h.storageDir, name)
}

func (h *handler) saveImage(image io.Reader, name string) (err error) {
	file, err := os.Create(h.pathToFile(name))
	if err != nil {
		return
	}
	defer file.Close()
	_, err = io.Copy(file, image)
	return
}

func (h *handler) lookupFile(name string) (file *os.File, err error) {
	return os.Open(h.pathToFile(name))
}

func getVars(r *http.Request) (vars map[string]interface{}, err error) {
	vars, ok := mux.Vars(r)
	if !ok {
		msg := fmt.Errorf(
			"did not find variables in request: %s", r.URL.String())
		log.Error(msg)
		err = msg
	}
	return
}
