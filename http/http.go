package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

const ShutdownTimeout = 1 * time.Second

// Server represents an HTTP server. It is meant to wrap all HTTP functionality
// used by the application so that dependent packages (such as cmd/wtfd) do not
// need to reference the "net/http" package at all.
type Server struct {
	ln     net.Listener
	server *http.Server
	//	router *mux.Router
	//	sc     *securecookie.SecureCookie

	// Bind address & domain for the server's listener.
	// If domain is specified, server is run on TLS using acme/autocert.
	Addr   string
	Domain string

	// Keys used for secure cookie encryption.
	HashKey  string
	BlockKey string

	// Servics used by the various HTTP routes.
	//	AuthService           wtf.AuthService
	//	DialService           wtf.DialService
	//	DialMembershipService wtf.DialMembershipService
	//	EventService          wtf.EventService
	//	UserService           wtf.UserService
}

func NewServer() *Server {

	s := &Server{}

	return s
}

func (s *Server) Open(port string) (err error) {
	if s.ln, err = net.Listen("tcp", port); err != nil {
		return err
	}
	// Begin serving requests on the listener. We use Serve() instead of
	// ListenAndServe() because it allows us to check for listen errors (such
	// as trying to use an already open port) synchronously.
	go s.server.Serve(s.ln)

	return nil
}

func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}
