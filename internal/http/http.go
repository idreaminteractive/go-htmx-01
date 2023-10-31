package http

import (
	"context"
	"encoding/gob"
	"fmt"
	"main/internal/config"
	"main/internal/db"
	"main/internal/services"
	"reflect"
	"strings"

	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"

	"github.com/gorilla/sessions"
)

const ShutdownTimeout = 1 * time.Second

type Server struct {
	echo                  *echo.Echo
	config                *config.EnvConfig
	sessionService        services.ISessionService
	authenticationService services.IAuthenticationService
	notesService          *services.NotesService
}
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

type EchoSetupStruct struct {
	SessionSecret string
	// default bool is false, so we generally want to enable it
	DisableCSRF bool
}

func handleSSE(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")

	// A function to send SSE messages to the client
	sendSSE := func(data string) {
		c.Response().Write([]byte("data: " + data + "\n\n"))
		c.Response().Flush()
	}

	// Simulate sending SSE messages (replace with your data source)
	go func() {
		for {
			sendSSE(fmt.Sprintf("<div>HTML @ %v</div>", time.Now().Format("20060102150405")))
			// You can replace this with actual data or events
			// Sleep for some time to simulate events
			time.Sleep(2 * time.Second)
		}
	}()

	// Ensure the connection remains open
	<-c.Request().Context().Done()
	return nil
}

func setupEcho(config EchoSetupStruct) *echo.Echo {
	// sets up echo with standard things
	// we attach it here in order to allow tests to use it as well.
	e := echo.New()
	// let's try sth different!
	// e.GET("/events", handleSSE)
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	// e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(
	//     rate.Limit(20),
	// )))

	gob.Register(services.SessionPayload{})
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(config.SessionSecret))))

	e.Use(middleware.Gzip())
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	e.Validator = &CustomValidator{validator: validate}
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

	ss := services.SessionService{SessionName: "_session", MaxAge: 3600}

	as := services.AuthenticationService{Queries: queries}
	// if we want to hide the queries?
	ns := services.InitNotesService(queries)

	// initialize the rest of our services
	s := &Server{
		authenticationService: &as,
		echo:                  e,
		sessionService:        &ss,
		config:                config,
		notesService:          ns,
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
	for _, item := range e.Router().Routes() {
		logrus.WithField("r", item).Info("")
	}

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
