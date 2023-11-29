package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Only group when necessary like w/ middlewares, etc.
func (s *Server) routes() {
	fs := http.FileServer(http.Dir("static"))
	s.router.Handle("/static/*", http.StripPrefix("/static/", fs))

	// health check routes
	// note, this is an example of using a closure for a route
	// to provide extra info, or repeatable routes
	s.router.Get("/healthz", s.handleAnyHealthz)

	s.router.NotFound(s.notFound)

	// Root routes
	s.router.Get("/", s.handleRootGet)
	s.router.Get("/login", s.handleLoginGet)
	s.router.Post("/login", s.handleLoginPost)

	s.router.Get("/logout", s.handleLogoutGet)

	s.router.Get("/register", s.handleRegisterGet)
	s.router.Post("/register", s.handleRegisterPost)

	s.router.Route("/chat", func(r chi.Router) {
		r.Use(s.requireAuthMiddleware)
		r.Get("/", s.handleChatGet)
		r.Get("/{id}", s.handleChatByIdGet)
		r.Post("/{id}", s.handleChatByIdPost)
		r.Post("/new", s.handleChatNewPost)

	})

	// add our events endpoint for sse
	s.router.Get("/events", s.services.SSEEventBus.ServeHTTP)

}
