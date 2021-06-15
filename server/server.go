// Package server contains all functions and methods for the main server that is
// the host to all API endpoints and static files for the chat.
package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sharpvik/log-go/v2"

	"github.com/sharpvik/ava/configs"
)

// Server is a wrapper around http.Server that allows us to define our own
// convenience methods.
type Server struct {
	server *http.Server
}

// NewServer returns appropriate server based on the mode and configs.
func NewServer(config *configs.Server) (s *Server) {
	s = NewBasicServer()
	s.setMode(config.DevMode)
	s.server.Handler = newServerHandler(config)
	return
}

// NewBasicServer returns a server with basic common-sense settings.
func NewBasicServer() *Server {
	return &Server{
		&http.Server{
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
	}
}

// Grace allows Server to shutdown gracefully based on an external done chan.
func (s *Server) Grace(done chan bool) {
	quit := make(chan os.Signal, 1)
	defer close(quit)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Debug("stopping server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s.server.SetKeepAlivesEnabled(false)

	if err := s.server.Shutdown(ctx); err != nil {
		log.Fatal("graceful server shutdown failed")
	}

	close(done)
}

// Serve starts the server.
func (s *Server) Serve() (err error) {
	log.Debugf("serving at %s ...", s.server.Addr)
	err = s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Errorf("server shut with error: %s", err)
	}
	return
}

// ServeWithGrace spawns the greceful shutfown monitor thread and then calls
// ListenAndServe on the server.
func (s *Server) ServeWithGrace(done chan bool) {
	go s.Grace(done)
	go s.Serve()
}
