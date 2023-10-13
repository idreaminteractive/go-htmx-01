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
	group.GET("/", s.handleDashboard)
	group.GET("/testing/", s.handleTest)
}
func (s *Server) handleTest(c echo.Context) error {
	component := views.Dashboard()
	base := views.Base(component)
	renderComponent(base, c)
	return nil
}

// will be the main page of the system
func (s *Server) handleDashboard(c echo.Context) error {
	component := views.Dashboard()
	base := views.Base(component)
	renderComponent(base, c)
	return nil
}
