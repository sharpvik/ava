package server

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/sharpvik/ava/auth"
	"github.com/sharpvik/log-go/v2"
	"github.com/sharpvik/mux"
)

// newServerHandler returns the main server handler responsible for the API.
func newServerHandler(apiKey string, storageDir http.Dir) http.Handler {
	handler := &handler{
		apiKey:     apiKey,
		storageDir: string(storageDir),
	}
	handler.authorized = authorizedHandler(handler)
	return handler
}

type handler struct {
	apiKey     string
	storageDir string
	authorized http.Handler
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if auth.Auth(h.apiKey, r) {
		h.authorized.ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}
}

func authorizedHandler(h *handler) http.Handler {
	rtr := mux.New().UseFunc(logRequest(""))

	rtr.Subrouter().
		Methods(http.MethodPost).
		Path("/upload/{ext:jpg|png}").
		HandleFunc(h.upload)

	rtr.Subrouter().
		Methods(http.MethodGet).
		Path(`/download/{name:.+\.(jpg|png)}`).
		HandleFunc(h.download)

	return rtr
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
