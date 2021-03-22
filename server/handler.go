package server

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/sharpvik/log-go/v2"
	"github.com/sharpvik/mux"
)

// newServerHandler returns the main server handler responsible for the API.
func newServerHandler(storageDir http.Dir) http.Handler {
	handler := &handler{string(storageDir)}

	rtr := mux.New().UseFunc(logRequest(""))

	rtr.Subrouter().
		Methods(http.MethodPost).
		Path("/upload/{ext:jpg|png}").
		HandleFunc(handler.upload)

	rtr.Subrouter().
		Methods(http.MethodGet).
		Path(`/download/{name:.+\.(jpg|png)}`).
		HandleFunc(handler.download)

	return rtr
}

type handler struct {
	storageDir string
}

func (h *handler) upload(w http.ResponseWriter, r *http.Request) {
	vars, err := getVars(r)
	if err != nil {
		return
	}

	name := newImageName(vars["ext"].(string))
	if err := h.saveImage(r.Body, name); err != nil {
		log.Errorf("failed to save image: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, name)
}

func (h *handler) download(w http.ResponseWriter, r *http.Request) {
	vars, err := getVars(r)
	if err != nil {
		return
	}

	name := vars["name"].(string)
	file, err := h.lookupFile(name)
	if err != nil {
		log.Errorf("failed to find/open image: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(name)))
	io.Copy(w, file)
}
