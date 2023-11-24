package http

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// all routes + structures are here
func (s *Server) handleSSE(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")

	// A function to send SSE messages to the client
	sendSSE := func(data string) {
		// can also send an event name here...
		if c.Response() != nil {
			_, err := c.Response().Write([]byte("data: " + data + "\n\n"))
			if err == nil {
				c.Response().Flush()
			}

		}

	}

	// eventually, this would be some link to listen to sth
	// for now, just check our messages
	go func() {

		for {
			// get our num messages
			count, err := s.services.ChatService.GetTotalMessagCount()
			if err != nil {
				logrus.Error(err)
				continue
			}

			sendSSE(fmt.Sprintf("%d", count))
			// You can replace this with actual data or events
			// Sleep for some time to simulate events
			time.Sleep(2 * time.Second)
		}
	}()

	// Ensure the connection remains open
	<-c.Request().Context().Done()
	return nil
}

// Only group when necessary like w/ middlewares, etc.
func (s *Server) routes() {

	// not sure how i feel about the handler names, but i mean, it's readable?
	s.echo.Static("/static", "static")
	s.echo.GET("/message-count", s.handleSSE)
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

	// Logged in routes
	chatGroup := s.echo.Group("/chat")
	chatGroup.Use(s.requireAuthMiddleware)
	chatGroup.GET("", s.handleChatGet)
	chatGroup.GET("/:id", s.handleChatByIdGet)
	chatGroup.POST("/:id", s.handleChatByIdPost)
	chatGroup.POST("/new", s.handleChatNewPost)

}
