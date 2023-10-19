package http

import (
	"main/internal/services"
	"main/internal/views"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) registerPublicRoutes() {
	s.echo.GET("/", s.handleHomeGet)
	s.echo.GET("/login/", s.handleLoginGet)
	s.echo.POST("/login/", s.handleLoginPost)

	s.echo.GET("/logout/", s.handleLogout)
}

func (s *Server) handleLogout(c echo.Context) error {
	// kill session + redirect (i should not need to post anywhere)
	// write a blank session
	s.sessionService.WriteSession(c, services.SessionPayload{})
	return c.Redirect(http.StatusMovedPermanently, "/")
}

// will be the main page of the system
// let's mirror our current live version that pulls in the stuff
func (s *Server) handleHomeGet(c echo.Context) error {

	renderComponent(views.Base(views.Home()), c)
	return nil
}

func (s *Server) handleLoginPost(c echo.Context) error {

	var user views.UserLoginDTO

	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(user); err != nil {

		// login failed, so let's send back bad request
		component := views.LoginForm(user, views.UserLoginFormErrors{Message: "Invalid login, please try again"})
		// return the view with our error

		renderComponent(component, c, 400)
		return nil
	}

	// create our user + id
	results, err := s.authenticationService.Authenticate(user)
	if err != nil {

		component := views.LoginForm(user, views.UserLoginFormErrors{Message: "Invalid login, please try again"})
		// return the view with our error
		renderComponent(component, c)
		return nil
	}

	// create our session + stuff
	s.sessionService.WriteSession(c, services.SessionPayload{UserId: int(results.ID), Email: user.Email})
	c.Response().Header().Set("HX-Redirect", "/dashboard")
	return c.NoContent(200)
}

func (s *Server) handleLoginGet(c echo.Context) error {
	// no errors or anything on initial bits.
	component := views.LoginPage(views.UserLoginDTO{}, views.UserLoginFormErrors{})
	base := views.Base(component)
	renderComponent(base, c)
	return nil
}
