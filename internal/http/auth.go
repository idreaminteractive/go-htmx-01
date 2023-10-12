package http

import (
	"main/internal/views"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (s *Server) registerAuthRoutes() {
	s.echo.GET("/login", s.handleLoginGet)
	s.echo.POST("/login", s.handleLoginPost)
}

func (s *Server) handleLoginPost(c echo.Context) error {
	var user views.UserLoginDTO
	err := c.Bind(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	logrus.WithField("user", user).Info("User")
	if err = c.Validate(user); err != nil {
		// should be an htmx post and we deal with that
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	// ok - it worked. so... what now?
	component := views.LoginForm(user, views.UserLoginFormErrors{"Invalid login, please try again"})
	templ.Handler(component).ServeHTTP(c.Response().Writer, c.Request())
	return nil
}

func (s *Server) handleLoginGet(c echo.Context) error {
	// no errors or anything on initial bits.
	component := views.LoginPage(views.UserLoginDTO{}, views.UserLoginFormErrors{})
	base := views.Base(component)
	templ.Handler(base).ServeHTTP(c.Response().Writer, c.Request())
	return nil
}
