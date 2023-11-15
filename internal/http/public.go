package http

import (
	"main/internal/services"
	"main/internal/views"
	"main/internal/views/dto"
	"net/http"

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
	if notes, err := s.services.NotesService.GetPublicNotes(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	} else {
		csrf_value := getCSRFValueFromContext(c)
		body := views.Home(views.HomePageData{Notes: notes})
		renderComponent(
			views.Base(
				views.BaseData{
					Body:  body,
					CSRF:  csrf_value,
					Title: "GoNotes",
				}),
			c)
	}

	return nil
}

func (s *Server) handleLoginPost(c echo.Context) error {

	var user dto.UserLoginDTO

	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := user.Validate(); err != nil {
		// login failed, so let's send back bad request
		// left off here - i feel like we can do better here
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

		renderComponent(component, c, 400)
		return nil
	}

	// create our user + id
	results, err := s.services.AuthenticationService.Authenticate(user)
	if err != nil {

		component := views.LoginPage(user, dto.UserLoginFormErrors{Message: "Invalid login, please try again"})
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

	component := views.LoginPage(dto.UserLoginDTO{}, dto.UserLoginFormErrors{})
	base := views.Base(views.BaseData{Body: component, CSRF: csrf_value, Title: "Login"})
	renderComponent(base, c)
	return nil
}
