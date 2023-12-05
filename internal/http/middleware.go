/*
# Route middlewares for site
*/
package http

import (
	"log/slog"
	"main/internal/helpers"
	"net/http"

	"github.com/go-chi/httplog/v2"
)

func (s *Server) requireAuthMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, err := s.services.SessionService.ReadSession(r)

		if err != nil {
			s.logger.Error("Error in getting session", err)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		if sess.UserId == 0 {
			s.logger.Error("Not logged in")
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// add our user id to logging
		httplog.LogEntrySetField(r.Context(), "user", slog.IntValue(sess.UserId))

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

func AddEnvMiddleware(cfg string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cr := helpers.ContextSave(r, "env", cfg)
			next.ServeHTTP(w, cr)
		})
	}

}
