package http

import (
	"main/internal/views"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (s *Server) requireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		logrus.Info("Check")
		sess, err := s.sessionService.ReadSession(c)
		if err != nil {
			logrus.WithField("err", err).Error("Error in getting session")
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		if sess.UserId == "" {
			logrus.Error("Not logged in")
			return c.Redirect(http.StatusMovedPermanently, "/login")
		}
		return next(c)
	}
}

func (s *Server) registerLoggedInRoutes(group *echo.Group) {
	group.GET("", s.handleDashboard)
}

// will be the main page of the system
// let's mirror our current live version that pulls in the stuff
func (s *Server) handleDashboard(c echo.Context) error {
	component := views.Hello("Dave")
	base := views.Base(component)
	renderComponent(base, c)
	return nil
}
