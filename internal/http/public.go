package http

import (
	"fmt"
	"main/internal/services"
	"main/internal/views"
	"main/internal/views/dto"
	"net/http"

	"github.com/a-h/templ"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sirupsen/logrus"

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
func (s *Server) handleRootGet(w http.ResponseWriter, r *http.Request) {

	// get our message count

	csrf_value := csrfFromRequest(r)
	count, err := s.services.ChatService.GetTotalMessagCount()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body := views.Root(count)

	// renderComponent(
	base := views.Base(
		views.BaseData{
			Body:  body,
			CSRF:  csrf_value,
			Title: "Chat App",
		},
	)
	base.Render(r.Context(), w)

}

func (s *Server) handleLoginPost(c echo.Context) error {

	var user dto.UserLoginDTO
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	var formErrors validation.Errors
	if err := user.Validate(); err != nil {
		formErrors = err.(validation.Errors)
		component := views.LoginScreen(views.LoginScreenProps{
			LastSubmission: user,
			Errors:         formErrors,
		})
		// JUst the tip
		c.Response().Header().Set("HX-Retarget", "#loginScreen")
		c.Response().Header().Set("HX-Reswap", "outerHTML")
		renderComponent(component, c)
		return nil
	}

	// create our user + id
	results, err := s.services.AuthenticationService.Authenticate(user)
	// auth fails
	if err != nil {
		component := views.LoginScreen(views.LoginScreenProps{
			LastSubmission: user,
			Errors: map[string]error{
				"email": validation.NewError("", "Invalid email or password"),
			},
		})
		c.Response().Header().Set("HX-Retarget", "#loginScreen")
		c.Response().Header().Set("HX-Reswap", "outerHTML")
		renderComponent(component, c)
		return nil
	}

	// successful login, lesgo.

	// create our session + stuff
	s.services.SessionService.WriteSession(c, services.SessionPayload{UserId: int(results.ID), Email: user.Email})

	c.Response().Header().Set("HX-Redirect", "/chat")

	return nil
}

func (s *Server) handleLoginGet(c echo.Context) error {
	// no errors or anything on initial bits.
	csrf_value := csrfFromRequest(c.Request())

	component := views.LoginScreen(views.LoginScreenProps{})
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
	component := views.RegisterForm(views.RegisterFormData{})
	base := views.Base(views.BaseData{Body: component, CSRF: csrfFromRequest(c.Request()), Title: "Register"})
	renderComponent(base, c)
	return nil
}

func (s *Server) handleRegisterPost(c echo.Context) error {
	var reg dto.RegisterDTO

	var formErrors validation.Errors

	if err := c.Bind(&reg); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := reg.Validate(); err != nil {
		formErrors = err.(validation.Errors)
	} else {
		// also check that password  === confirm password
		fmt.Printf("%q - %q\n", reg.Password, reg.ConfirmPassword)
		if reg.Password != reg.ConfirmPassword {
			// err.
			formErrors = map[string]error{
				"confirm_password": validation.NewError("", "The confirmed password must match your password"),
			}
		}
	}

	// let's hash our password + then check and see if the user already exists (we hash first to prevent timing attacks)
	user, err := s.services.AuthenticationService.Register(reg)
	if err != nil {
		logrus.Errorf("%s", err)
		// actually check if it exists or not
		formErrors = map[string]error{
			"email": validation.NewError("", "That user already exists "),
		}
	}

	if formErrors.Filter() != nil {
		component := views.RegisterForm(views.RegisterFormData{Previous: reg, Errors: formErrors})
		c.Response().Header().Set("HX-Retarget", "#registerForm")
		c.Response().Header().Set("HX-Reswap", "outerHTML")
		renderComponent(component, c)
		return nil

	}
	// ok - return success!
	logrus.Info("Successful registration!")
	s.services.SessionService.WriteSession(c, services.SessionPayload{UserId: int(user.ID), Email: user.Email})

	c.Response().Header().Set("HX-Redirect", "/chat")
	return nil
}
