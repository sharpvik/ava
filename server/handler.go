package server

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/sharpvik/log-go/v2"
	"github.com/sharpvik/mux"

	"github.com/sharpvik/ava/auth"
	"github.com/sharpvik/ava/configs"
)

type handler struct {
	*configs.Server
	authorized http.Handler
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !auth.Auth(h.APIKey, r) {
		respondWithStatusAndMessage(w,
			http.StatusUnauthorized, "401 Unauthorized - Wrong API Key")
		return
	}
	if !h.bodySizeUnderLimit(r) {
		respondWithStatusAndMessage(w,
			http.StatusBadRequest, "400 Bad Request - Image Size Too Big")
		return
	}
	h.authorized.ServeHTTP(w, r)
}

// newServerHandler returns the main server handler responsible for the API.
func newServerHandler(config *configs.Server) http.Handler {
	handler := &handler{
		Server: config,
	}
	handler.authorized = handler.authorizedHandler()
	return handler
}

func (h *handler) authorizedHandler() http.Handler {
	rtr := mux.New().UseFunc(logRequest(""))

	rtr.Subrouter().
		Methods(http.MethodPost).
		Path(`/upload/{ext:\w+}`).
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

	ext := vars["ext"].(string)
	if !h.extensionIsAllowed(ext) {
		respondWithStatusAndMessage(w,
			http.StatusForbidden, "403 Forbidden - Extension Not Allowed")
		return
	}

	name := newImageName(ext)
	if err := h.saveImage(r.Body, name); err != nil {
		log.Errorf("failed to save image: %s", err)
		respondWithStatusAndMessage(w,
			http.StatusInternalServerError,
			"500 Internal Server Error - Failed To Save Image")
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
		respondWithStatusAndMessage(w,
			http.StatusInternalServerError,
			"500 Internal Server Error - Cannot Find Image")
		return
	}

	defer file.Close()
	w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(name)))
	io.Copy(w, file)
}
