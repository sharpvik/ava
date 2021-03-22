package server

import "github.com/sharpvik/log-go/v2"

// setMode takes a boolean flag called devMode and uses it to decide on the
// proper server mode.
func (s *Server) setMode(devMode bool) {
	var mode string
	if devMode {
		mode = s.toDevServer()
	} else {
		mode = s.toProdServer()
	}
	log.Debugf("server running in %s mode ...", mode)
}

// toDevServer returns a server that runs on localhost and is unavailable from
// the outside.
func (s *Server) toDevServer() string {
	s.server.Addr = "127.0.0.1:42069"
	return "development"
}

// toProdServer returns a server that will be listening on the publically
// available port 4321.
func (s *Server) toProdServer() string {
	s.server.Addr = "0.0.0.0:42069"
	return "production"
}
