/*
# Route middlewares for site
*/
package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

func (s *Server) requireAuthMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, err := s.services.SessionService.ReadSession(r)

		if err != nil {
			logrus.WithField("err", err).Error("Error in getting session")
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		if sess.UserId == 0 {
			logrus.Error("Not logged in")
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		ctx := context.WithValue(r.Context(), "userId", sess.UserId)
		next.ServeHTTP(w, r.WithContext(ctx))
		return
	})
}

func (s *Server) getUserIdFromCTX(r *http.Request) int {

	userId, err := contextGet(r, "userId")
	if err != nil {
		return 0
	}
	return userId.(int)
}

func contextGet(r *http.Request, key string) (interface{}, error) {
	val := r.Context().Value(key)
	if val == nil {
		return nil, fmt.Errorf("no value exists in the context for key %q", key)
	}
	return val, nil
}

func contextSave(r *http.Request, key string, val interface{}) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, key, val) // nolint:staticcheck
	return r.WithContext(ctx)
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
