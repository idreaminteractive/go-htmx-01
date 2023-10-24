package http

import (
	"main/internal/views"

	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (s *Server) requireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		sess, err := s.sessionService.ReadSession(c)
		if err != nil {
			logrus.WithField("err", err).Error("Error in getting session")
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}

		if sess.UserId == 0 {
			logrus.Error("Not logged in")
			return c.Redirect(http.StatusFound, "/login")
		}

		// for _, val := range c.Echo().Router().Routes() {
		// 	logrus.WithField("val", val).Info("I should goto thing?")
		// }

		return next(c)
	}
}

func (s *Server) registerLoggedInRoutes(group *echo.Group) {
	group.GET("", s.handleDashboard)
	noteGroup := group.Group("/note")

	s.registerNoteRoutes(noteGroup)

}

// will be the main page of the system
func (s *Server) handleDashboard(c echo.Context) error {

	// find our logged in user to get their personal notes
	sp, err := s.sessionService.ReadSession(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not read session")

	}

	userNotes, err := s.notesService.GetNotesForUserId(sp.UserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not fetch notes for user")
	}
	csrf_value := getCSRFValueFromContext(c)
	component := views.Dashboard(userNotes)
	base := views.Base(component, csrf_value)
	renderComponent(base, c)
	return nil
}
