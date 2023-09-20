package http

import (
	"main/views/routes"

	"github.com/labstack/echo/v4"
)

func (s *Server) registerHomeRoutes() {

	s.echo.GET("/", s.handleHome)
}

func (s *Server) handleHome(c echo.Context) error {
	return routes.HomePage().Render(c.Response().Writer)

}
