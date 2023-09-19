package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) registerHomeRoutes() {

	s.echo.GET("/", s.handleHome)
}

func (s *Server) handleHome(c echo.Context) error {
	return c.String(http.StatusOK, "Hello world!")

}
