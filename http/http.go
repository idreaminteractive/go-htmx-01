package http

import (
	"context"
	// "net"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

const ShutdownTimeout = 1 * time.Second

// Server represents an HTTP server. It is meant to wrap all HTTP functionality
// used by the application so that dependent packages (such as cmd/wtfd) do not
// need to reference the "net/http" package at all.
type Server struct {
	server *http.Server
	router *chi.Mux

	Potato string
	// router *mux.Router
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

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	s := &Server{
		router: r,
		server: &http.Server{},
	}

	s.Potato = "I am a potaot"

	// handlers + routers
	// all routers + sub routers come off of the server struct
	// we can organize w/ files and simple extend the struct
	// (https://github.com/benbjohnson/wtf/blob/main/http/dial.go#L210 for example)
	r.Get("/", s.handleRoot)

	return s
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(s.Potato))
}

func (s *Server) Open(port string) (err error) {

	http.ListenAndServe(port, s.router)
	return nil
}

func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}
