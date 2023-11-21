package http

// all routes + structures are here

// Only group when necessary like w/ middlewares, etc.
func (s *Server) routes() {

	// not sure how i feel about the handler names, but i mean, it's readable?
	s.echo.Static("/static", "static")

	// health check routes
	// note, this is an example of using a closure for a route
	// to provide extra info, or repeatable routes
	s.echo.Any("/healthz", s.handleAnyHealthz())

	// Root routes
	s.echo.GET("/", s.handleRootGet)
	s.echo.GET("/login", s.handleLoginGet)
	s.echo.POST("/login", s.handleLoginPost)

	s.echo.GET("/logout", s.handleLogoutGet)

	s.echo.GET("/register", s.handleRegisterGet)
	s.echo.POST("/register", s.handleRegisterPost)

	// once we're logged in, where to?
	s.echo.GET("/chat", s.handleChatGet)

	// s.echo.GET("/message-count", s.handleMessageCountGet)

	// Authenticated routes
	// authGroup := s.echo.Group("/dashboard")
	// authGroup.Use(s.requireAuthMiddleware)
	// authGroup.GET("", s.handleDashboardGet)

	// authGroup.POST("/note", s.handleNotePost)
	// authGroup.DELETE("/note/:id", s.handleNoteByIdDelete)
	// authGroup.GET("/note/:id/edit", s.handleNoteByIdEditGet)
	// authGroup.PUT("/note/:id", s.handleNoteByIdPut)

}
