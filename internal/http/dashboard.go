package http

import (
	"main/internal/views"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func (s *Server) registerLoggedInRoutes() {
	s.echo.GET("/dashboard", s.handleDashboard)

}

// will be the main page of the system
// let's mirror our current live version that pulls in the stuff
func (s *Server) handleDashboard(c echo.Context) error {
	component := views.Hello("Dave")
	base := views.Base(component)
	templ.Handler(base).ServeHTTP(c.Response().Writer, c.Request())
	return nil
}
