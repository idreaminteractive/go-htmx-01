package http

import (
	"context"
	"encoding/gob"
	"main/internal/config"
	"main/internal/db"
	"main/internal/services"

	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/gorilla/sessions"
)

const ShutdownTimeout = 1 * time.Second

type Server struct {
	echo                  *echo.Echo
	config                *config.EnvConfig
	sessionService        services.ISessionService
	authenticationService services.IAuthenticationService
}
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func setupEcho(sessionSecret string) *echo.Echo {
	// sets up echo with standard things
	// we attach it here in order to allow tests to use it as well.
	e := echo.New()
	e.Pre(middleware.AddTrailingSlash())
	gob.Register(services.SessionPayload{})
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(sessionSecret))))

	e.Use(middleware.Gzip())
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "form:csrf",
	}))
	return e
}

func NewServer(config *config.EnvConfig, queries *db.Queries) *Server {
	// This is where we initialize all our services and attach to our
	// server
	e := setupEcho(config.SessionSecret)

	ss := services.SessionService{SessionName: "_session", MaxAge: 3600}

	as := services.AuthenticationService{Queries: queries}

	// initialize the rest of our services
	s := &Server{
		authenticationService: &as,
		echo:                  e,
		sessionService:        &ss,
		config:                config,
	}

	// for now, this is fine - we'll set some monster caching later on
	e.Static("/static", "static")

	// health check routes
	e.HEAD("/_health/", s.healthCheckRoute)
	e.GET("/_health/", s.healthCheckRoute)

	s.registerPublicRoutes()

	loggedInGroup := e.Group("/dashboard")
	loggedInGroup.Use(s.requireAuth)

	s.registerLoggedInRoutes(loggedInGroup)
	return s
}
func (s *Server) healthCheckRoute(c echo.Context) error {

	return c.String(http.StatusOK, "ok")

}

func (s *Server) Open(port string) (err error) {

	s.echo.Logger.Fatal(s.echo.Start(port))

	return nil

}

func (s *Server) Close() error {

	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)

	defer cancel()

	return s.echo.Shutdown(ctx)

}
