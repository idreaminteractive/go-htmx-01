package http

import (
	// "context"
	"context"
	"main/config"

	// "net"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const ShutdownTimeout = 1 * time.Second

// Server represents an HTTP server. It is meant to wrap all HTTP functionality
// used by the application so that dependent packages (such as cmd/wtfd) do not
// need to reference the "net/http" package at all.
type Server struct {
	echo   *echo.Echo
	config *config.EnvConfig
}
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
func NewServer(config *config.EnvConfig) *Server {
	// This is where we initialize all our services and attach to our
	// server

	e := echo.New()

	s := &Server{
		echo:   e,
		config: config,
	}

	// for now, this is fine - we'll set some monster caching later on
	e.Static("/static", "static")
	e.Use(middleware.Gzip())
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())

	e.Use(middleware.Recover())

	e.HEAD("/_health", s.healthCheckRoute)

	e.GET("/_health", s.healthCheckRoute)

	s.registerAuthRoutes()
	s.registerRootRoutes()
	return s
}
func (s *Server) healthCheckRoute(c echo.Context) error {

	return c.String(http.StatusOK, "ok")

}

func (s *Server) Open(port string) (err error) {

	s.echo.Logger.Fatal(s.echo.Start(port))

	return nil

}

func (s *Server) Close() error {

	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)

	defer cancel()

	return s.echo.Shutdown(ctx)

}
