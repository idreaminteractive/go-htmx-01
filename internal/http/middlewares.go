/*
# Route middlewares for site
*/
package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (s *Server) requireAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		sess, err := s.services.SessionService.ReadSession(c)
		if err != nil {
			logrus.WithField("err", err).Error("Error in getting session")
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}

		if sess.UserId == 0 {
			logrus.Error("Not logged in")
			return c.Redirect(http.StatusFound, "/login")
		}

		return next(c)
	}
}
