/*
# Route middlewares for site
*/
package http

import (
	"main/internal/helpers"
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

		cr := helpers.ContextSave(r, "userId", sess.UserId)
		next.ServeHTTP(w, cr)
	})
}

func (s *Server) getUserIdFromCTX(r *http.Request) int {

	userId, err := helpers.ContextGet(r, "userId")
	if err != nil {
		return 0
	}
	return userId.(int)
}
