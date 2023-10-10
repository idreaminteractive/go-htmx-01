package http

import (
	"main/views"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func (s *Server) registerRootRoutes() {
	s.echo.GET("/", s.handleRootGet)

}

// will be the main page of the system
// let's mirror our current live version that pulls in the stuff
func (s *Server) handleRootGet(c echo.Context) error {
	component := views.Login()
	base := views.Base(component)
	templ.Handler(base).ServeHTTP(c.Response().Writer, c.Request())
	return nil
}
