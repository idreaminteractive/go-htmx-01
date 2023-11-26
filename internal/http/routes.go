package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// all routes + structures are here
func (s *Server) handleSSE(c echo.Context) error {

	flusher, ok := c.Response().Writer.(http.Flusher)

	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Streaming unsupported!")

	}
	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	// A function to send SSE messages to the client
	// sendSSE := func(data string) {
	// 	// can also send an event name here...
	// 	if c.Response() != nil {
	// 		_, err := c.Response().Write([]byte("data: " + data + "\n\n"))
	// 		if err == nil {
	// 			c.Response().Flush()
	// 		}

	// 	}

	// }

	// // eventually, this would be some link to listen to sth
	// // for now, just check our messages
	// go func() {

	// 	for {
	// 		// get our num messages
	// 		count, err := s.services.ChatService.GetTotalMessagCount()
	// 		if err != nil {
	// 			logrus.Error(err)
	// 			continue
	// 		}

	// 		sendSSE(fmt.Sprintf("%d", count))
	// 		// You can replace this with actual data or events
	// 		// Sleep for some time to simulate events
	// 		time.Sleep(2 * time.Second)
	// 	}
	// }()

	// // Ensure the connection remains open
	// <-c.Request().Context().Done()
	// register connection message channel

	// signal the Server of a new client connection
	// sseServer.NewClientsChannel <- messageChan

	// keeping the connection alive with keep-alive protocol
	keepAliveTicker := time.NewTicker(15 * time.Second)
	keepAliveMsg := ":keepalive\n"
	notify := c.Request().Context().Done()

	// listen to signal to close and unregister
	go func() {
		<-notify
		// sseServer.ClosingClientsChannel <- messageChan
		keepAliveTicker.Stop()
	}()

	defer func() {
		// sseServer.ClosingClientsChannel <- messageChan
	}()

	for {
		select {
		//receiving a message from the Kafka channel.
		case messageCountEvent := <-s.services.ChatService.MessageChannel:
			// Write to the ResponseWriter in SSE compatible format
			c.Response().Write([]byte(fmt.Sprintf("data: %s\n\n", messageCountEvent)))
			c.Response().Flush()
		case <-keepAliveTicker.C:
			fmt.Fprintf(c.Response().Writer, keepAliveMsg)
			flusher.Flush()
		}

	}

	// return nil
}

// Only group when necessary like w/ middlewares, etc.
func (s *Server) routes() {
	fs := http.FileServer(http.Dir("static"))
	s.router.Handle("/static/*", http.StripPrefix("/static/", fs))
	// not sure how i feel about the handler names, but i mean, it's readable?

	// s.router.Get("/message-count", s.handleSSE)
	// health check routes
	// note, this is an example of using a closure for a route
	// to provide extra info, or repeatable routes
	s.router.Get("/healthz", s.handleAnyHealthz)

	// Root routes
	s.router.Get("/", s.handleRootGet)
	// s.echo.GET("/login", s.handleLoginGet)
	// s.echo.POST("/login", s.handleLoginPost)

	// s.echo.GET("/logout", s.handleLogoutGet)

	// s.echo.GET("/register", s.handleRegisterGet)
	// s.echo.POST("/register", s.handleRegisterPost)

	// // Logged in routes
	// chatGroup := s.echo.Group("/chat")
	// chatGroup.Use(s.requireAuthMiddleware)
	// chatGroup.GET("", s.handleChatGet)
	// chatGroup.GET("/:id", s.handleChatByIdGet)
	// chatGroup.POST("/:id", s.handleChatByIdPost)
	// chatGroup.POST("/new", s.handleChatNewPost)

}
