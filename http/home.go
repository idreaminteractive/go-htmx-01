package http

import (
	"github.com/labstack/echo/v4"
	"main/views/routes"
)

func (s *Server) registerHomeRoutes() {

	s.echo.GET("/", s.handleHome)
}

func (s *Server) handleHome(c echo.Context) error {
	return routes.Home().Render(c.Response().Writer)
	//return c.String(http.StatusOK, "Hello world!")
	//
}
