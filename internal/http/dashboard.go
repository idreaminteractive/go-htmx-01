package http

import (
	"main/internal/views"

	"github.com/labstack/echo/v4"
)

// will be the main page of the system
func (s *Server) handleDashboardGet(c echo.Context) error {

	csrf_value := csrfFromRequest(c.Request())
	component := views.Dashboard()
	base := views.Base(views.BaseData{Body: component, CSRF: csrf_value})
	renderComponent(base, c)
	return nil
}
