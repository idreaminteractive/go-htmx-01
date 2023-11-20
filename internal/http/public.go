package http

import (
	"main/internal/services"
	"main/internal/views"
	"main/internal/views/dto"
	"net/http"

	"github.com/a-h/templ"
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/labstack/echo/v4"
)

func (s *Server) handleLogoutGet(c echo.Context) error {
	// kill session + redirect (i should not need to post anywhere)
	// write a blank session
	s.services.SessionService.WriteSession(c, services.SessionPayload{})
	return c.Redirect(http.StatusMovedPermanently, "/")
}

// will be the main page of the system
// let's mirror our current live version that pulls in the stuff
func (s *Server) handleRootGet(c echo.Context) error {

	// get our public notes
	// if notes, err := s.services.NotesService.GetPublicNotes(); err != nil {
	// 	return echo.NewHTTPError(http.StatusInternalServerError, err)
	// } else {
	// 	csrf_value := getCSRFValueFromContext(c)
	// 	body := views.Home(views.HomePageData{Notes: notes})
	// 	renderComponent(
	// 		views.Base(
	// 			views.BaseData{
	// 				Body:  body,
	// 				CSRF:  csrf_value,
	// 				Title: "GoNotes",
	// 			}),
	// 		c)
	// }
	csrf_value := getCSRFValueFromContext(c)
	body := views.Root()

	renderComponent(
		views.Base(
			views.BaseData{
				Body:  body,
				CSRF:  csrf_value,
				Title: "GoNotes",
			},
		),
		c)

	return nil
}

func (s *Server) handleLoginPost(c echo.Context) error {

	var user dto.UserLoginDTO

	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := user.Validate(); err != nil {

		// validation failed.
		component := views.Base(views.BaseData{
			Body: views.LoginPage(
				views.LoginPageData{
					LoginForm: views.LoginForm(views.LoginFormData{
						Errors:   err.(validation.Errors),
						Defaults: user,
					}),
				}),
		})
		// return the view with our error
		// note - it's a 200 message ALWAYS
		renderComponent(component, c)
		return nil
	}

	// create our user + id
	results, err := s.services.AuthenticationService.Authenticate(user)
	// auth fails
	if err != nil {
		// not a fan of this.... can prob clean it up
		component := views.Base(views.BaseData{
			Body: views.LoginPage(
				views.LoginPageData{
					LoginForm: views.LoginForm(views.LoginFormData{
						Errors: map[string]error{
							"email": validation.NewError("", "Email or password is invalid "),
						},
						Defaults: user,
					}),
				}),
		})
		// return the view with our error
		renderComponent(component, c)
		return nil
	}

	// create our session + stuff
	s.services.SessionService.WriteSession(c, services.SessionPayload{UserId: int(results.ID), Email: user.Email})

	c.Response().Header().Set("HX-Redirect", "/dashboard")

	return c.NoContent(http.StatusOK)
}

func (s *Server) handleLoginGet(c echo.Context) error {
	// no errors or anything on initial bits.
	csrf_value := getCSRFValueFromContext(c)

	component := views.LoginScreen()
	base := views.Base(views.BaseData{Body: component, CSRF: csrf_value, Title: "Login"})
	renderComponent(base, c)
	return nil
}

func (s *Server) handleMessageCountGet(c echo.Context) error {

	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().Header().Set("Content-Type", "text/event-stream")
	component := views.MessageCount(13)
	c.Response().Writer.WriteHeader(200)
	templ.Handler(component).ServeHTTP(c.Response().Writer, c.Request())
	// renderComponent(component, c)
	return nil
}

func (s *Server) handleRegisterGet(c echo.Context) error {
	// no errors or anything on initial bits.
	component := views.RegisterForm()
	base := views.Base(views.BaseData{Body: component, CSRF: getCSRFValueFromContext(c), Title: "Register"})
	renderComponent(base, c)
	return nil
}

func (s *Server) handleRegisterPost(c echo.Context) error {
	// no errors or anything on initial bits.

	return echo.NewHTTPError(http.StatusNotFound)
}
