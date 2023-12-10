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
	s.router.Get("/healthz", s.handleAnyHealthz)

	s.router.NotFound(s.notFound)

	// Root routes
	s.router.Get("/", s.handleRootGet)
	s.router.Get("/login", s.handleLoginGet)
	s.router.Post("/login", s.handleLoginPost)

	s.router.Get("/logout", s.handleLogoutGet)

	s.router.Get("/register", s.handleRegisterGet)
	s.router.Post("/register", s.handleRegisterPost)

// authenticated route
s.router.Route("/loggedin", func(r chi.Router) {
	r.Use(s.requireAuthMiddleware)
	r.Get("/", s.handleLoggedInGet)

})


	// create our ws stuff
	// hub := ws.NewHub()
	// go hub.Run()
	// // this could be under a middleware too so the ws can get our user.
	// s.router.Get("/chatws", func(w http.ResponseWriter, r *http.Request) {
	// 	sess, err := s.services.SessionService.ReadSession(r)

	// 	if err != nil {
	// 		s.logger.Error("Error in getting session", err)
	// 		http.Redirect(w, r, "/login", http.StatusFound)
	// 		return
	// 	}
	// 	spew.Dump(sess)
	// 	ws.ServeWs(hub, w, r)
	// })

	s.router.Route("/chat", func(r chi.Router) {
		r.Use(s.requireAuthMiddleware)
		r.Get("/", s.handleChatGet)
		r.Get("/{id}", s.handleChatByIdGet)
		r.Post("/{id}", s.handleChatByIdPost)
		r.Post("/new", s.handleChatNewPost)

	})

	// add our events endpoint for sse
	s.router.Get("/hmr", s.services.SSEEventBus.ServeHTTP)

}
