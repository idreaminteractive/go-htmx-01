package http

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"reflect"
	"strings"

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

const tagCustom = "errmsg"

func errorTagFunc[T interface{}](obj interface{}, snp string, fieldname, actualTag string) error {
	o := obj.(T)
	fmt.Printf("o - %v\n", o)
	if !strings.Contains(snp, fieldname) {
		fmt.Printf("Does not conainer, %v, %v\n", snp, fieldname)
		return nil
	}

	fieldArr := strings.Split(snp, ".")
	rsf := reflect.TypeOf(o)

	for i := 1; i < len(fieldArr); i++ {
		field, found := rsf.FieldByName(fieldArr[i])
		if found {
			if fieldArr[i] == fieldname {
				customMessage := field.Tag.Get(tagCustom)
				if customMessage != "" {
					return fmt.Errorf("%s: %s (%s)", fieldname, customMessage, actualTag)
				}
				return nil
			} else {
				if field.Type.Kind() == reflect.Ptr {
					// If the field type is a pointer, dereference it
					rsf = field.Type.Elem()
				} else {
					rsf = field.Type
				}
			}
		}
	}
	return nil
}

func ValidateFunc[T interface{}](obj interface{}, validate *validator.Validate) (errs error) {
	o := obj.(T)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in Validate:", r)
			errs = fmt.Errorf("can't validate %+v", r)
		}
	}()

	if err := validate.Struct(o); err != nil {
		errorValid := err.(validator.ValidationErrors)
		for _, e := range errorValid {
			// snp  X.Y.Z
			snp := e.StructNamespace()
			fmt.Printf("\n%v - %v - %v - %v\n", obj, snp, e.Field(), e.ActualTag())
			errmgs := errorTagFunc[T](obj, snp, e.Field(), e.ActualTag())
			if errmgs != nil {
				errs = errors.Join(errs, fmt.Errorf("%w", errmgs))
			} else {
				errs = errors.Join(errs, fmt.Errorf("%w", e))
			}
		}
	}

	if errs != nil {
		return errs
	}

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

	note := dto.CreateNoteDTO{Content: "", IsPublic: "on"}
	// set that we're looking for form

	errs := note.Validate()
	validation.ErrorTag = "form"
	fmt.Printf("ERRS1: %+v\n", errs)
	errs = note.Validate()
	// fmt.Printf("")
	// errs := ValidateFunc[dto.CreateNoteDTO](note, validate)
	fmt.Printf("ERRS2: %+v\n", errs)
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
