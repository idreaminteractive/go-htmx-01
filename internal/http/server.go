package http

import (
	"context"
	"encoding/gob"

	"main/internal/config"
	"main/internal/db"
	"main/internal/services"

	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/gorilla/sessions"
)

const ShutdownTimeout = 1 * time.Second

// i want to
type Server struct {
	echo   *echo.Echo
	config *config.EnvConfig

	services *services.ServiceLocator
}

type EchoSetupStruct struct {
	SessionSecret string
	// default bool is false, so we generally want to enable it
	DisableCSRF bool
}

func setupEcho(config EchoSetupStruct) *echo.Echo {
	// sets up echo with standard things
	// we attach it here in order to allow tests to use it as well.
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	// e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(
	//     rate.Limit(20),
	// )))

	gob.Register(services.SessionPayload{})
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(config.SessionSecret))))

	e.Use(middleware.Gzip())

	// set that we're looking for form
	validation.ErrorTag = "form"
	// errs := note.Validate()

	// if errs != nil {
	// 	terr := errs.(validation.Errors)
	// 	fmt.Println(terr["content"])
	// 	fmt.Println("------")
	// 	fmt.Println(terr["potato"])
	// 	fmt.Println("------")

	// }

	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	if !config.DisableCSRF {
		e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
			TokenLookup: "header:X-CSRFToken",
			// X-CSRFToken
		}))

	}

	return e
}

func NewServer(config *config.EnvConfig, queries *db.Queries) *Server {
	// This is where we initialize all our services and attach to our
	// server
	e := setupEcho(EchoSetupStruct{SessionSecret: config.SessionSecret})

	// setup our service locator
	sl := services.ServiceLocator{}

	sl.AuthenticationService = services.InitAuthService(&sl, queries)
	sl.SessionService = services.InitSessionService(&sl, "_session", 3600)
	sl.NotesService = services.InitNotesService(&sl, queries)

	// initialize the rest of our services
	s := &Server{
		echo:     e,
		config:   config,
		services: &sl,
	}

	// for now, this is fine - we'll set some monster caching later on
	e.Static("/static", "static")

	// health check routes
	e.HEAD("/_health", s.healthCheckRoute)
	e.GET("/_health", s.healthCheckRoute)

	s.registerPublicRoutes()

	loggedInGroup := e.Group("/dashboard")
	loggedInGroup.Use(s.requireAuth)

	s.registerLoggedInRoutes(loggedInGroup)

	// print the routes
	// for _, item := range e.Router().Routes() {
	// 	logrus.WithField("r", item).Info("")
	// }

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

// safe csrf getting
func getCSRFValueFromContext(c echo.Context) string {
	context := c.Get(middleware.DefaultCSRFConfig.ContextKey)
	if context == nil {
		// we don't have anything here, use blank string
		return ""
	}
	return context.(string)
}
