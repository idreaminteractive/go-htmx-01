package http

import (
	"context"
	"encoding/gob"
	"fmt"
	"reflect"

	"main/internal/config"
	"main/internal/db"
	"main/internal/services"
	"main/internal/views/dto"

	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
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
	notesService          *services.NotesService
}
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	// this is where we can write our custom overlap
	return cv.validator.Struct(i)
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

	note := dto.CreateNoteDTO{Content: "", IsPublic: "on"}
	// set that we're looking for form
	validation.ErrorTag = "form"
	errs := note.Validate()
	fmt.Println(reflect.TypeOf(errs))
	if errs != nil {
		// cast to validation errs + make sure it's good
		terr := errs.(validation.Errors)
		fmt.Println(terr["content"])
		fmt.Println("------")
		fmt.Println(terr["potato"])
		fmt.Println("------")
		// ok - les go + pass this in?
		// for key, errObject := range terr {
		// 	fmt.Println(key)
		// 	fmt.Printf("%+v", errObject.(validation.Error).Error())
		// }

	}

	// validate.RegisterTagNameFunc(func(fld reflect.StructField) string {

	// 	fmt.Printf("fld: %v\n", fld)
	// 	fmt.Printf("tag: %v\n", fld.Tag)
	// 	name := strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
	// 	if name == "-" {
	// 		return ""
	// 	}
	// 	return name
	// })

	e.Validator = &CustomValidator{validator: validate}

	// test out our validator systems.

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
