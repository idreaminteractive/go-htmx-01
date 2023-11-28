/*
# Route middlewares for site
*/
package http

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
)

func (s *Server) requireAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, err := s.services.SessionService.ReadSession(r)
		if err != nil {
			logrus.WithField("err", err).Error("Error in getting session")

			return
		}

		if sess.UserId == 0 {
			logrus.Error("Not logged in")
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		// todo - add id to ctx

		context.WithValue(r.Context(), "userId", sess.UserId)

		next.ServeHTTP(w, r)
		return
	})
}

//   (next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {

// 		sess, err := s.services.SessionService.ReadSession(r)
// 		if err != nil {
// 			logrus.WithField("err", err).Error("Error in getting session")
// 			return echo.NewHTTPError(http.StatusUnauthorized, err)
// 		}

// 		if sess.UserId == 0 {
// 			logrus.Error("Not logged in")
// 			return c.Redirect(http.StatusFound, "/login")
// 		}
// 		// todo - add id to ctx
// 		c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), "userId", sess.UserId)))

// 		return next(c)
// 	}
// }
