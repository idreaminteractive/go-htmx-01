package http

import (
	"main/internal/views"

	"net/http"

	"github.com/labstack/echo/v4"
)

// will be the main page of the system
func (s *Server) handleDashboardGet(c echo.Context) error {

	// find our logged in user to get their personal notes
	_, err := s.services.SessionService.ReadSession(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not read session")

	}

	csrf_value := getCSRFValueFromContext(c)
	component := views.Dashboard()
	base := views.Base(views.BaseData{Body: component, CSRF: csrf_value})
	renderComponent(base, c)
	return nil
}
