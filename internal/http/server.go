package http

import (
	"context"
	"encoding/gob"
	"fmt"

	"main/internal/config"
	"main/internal/db"
	"main/internal/services"

	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"

	"github.com/gorilla/csrf"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const ShutdownTimeout = 1 * time.Second

// i want to
type Server struct {
	router   *chi.Mux
	server   *http.Server
	config   *config.EnvConfig
	services *services.ServiceLocator
}

type ServerSetupStruct struct {
	SessionSecret string
	// default bool is false, so we generally want to enable it
	DisableCSRF bool
}

func setupServer(config ServerSetupStruct) *chi.Mux {
	// sets up echo with standard things
	// we attach it here in order to allow tests to use it as well.
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(60 * time.Second))
	gob.Register(services.SessionPayload{})
	// e.Use(session.Middleware(sessions.NewCookieStore([]byte(config.SessionSecret))))
	r.Use(middleware.Compress(5))

	validation.ErrorTag = "form"
	if !config.DisableCSRF {
		csrfMiddleware :=
			csrf.Protect([]byte("32-byte-long-auth-key"))
		r.Use(csrfMiddleware)

		// e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		// 	TokenLookup: "header:X-CSRFToken",
		// }))

	}

	return r
}

func NewServer(config *config.EnvConfig, queries *db.Queries) *Server {
	// This is where we initialize all our services and attach to our
	// server
	r := setupServer(ServerSetupStruct{SessionSecret: config.SessionSecret})

	// setup our service locator
	sl := services.ServiceLocator{}

	sl.AuthenticationService = services.InitAuthService(&sl, queries)
	sl.SessionService = services.InitSessionService(&sl, "_session", 3600)
	sl.ChatService = services.InitChatService(&sl, queries)
	// initialize the rest of our services + http server
	s := &Server{
		router:   r,
		config:   config,
		services: &sl,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%s", config.GoPort),
			Handler: r,
		},
	}
	s.routes()

	return s
}

// example of the route closures
func (s *Server) handleAnyHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))

}

func (s *Server) Open(port string) (err error) {

	return s.server.ListenAndServe()

}

func (s *Server) Close() error {

	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)

	defer cancel()
	return s.server.Shutdown(ctx)

}

// safe csrf getting
func getCSRFValueFromContext(c echo.Context) string {
	// context := c.Get(middleware.ContextKey)
	// if context == nil {
	// 	// we don't have anything here, use blank string
	// 	return ""
	// }
	// return context.(string)
	return ""
}
