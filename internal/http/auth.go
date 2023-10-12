package http

import (
	// "github.com/a-h/templ"
	"main/internal/views"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type UserLoginDTO struct {
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required"`
}

func (s *Server) registerAuthRoutes() {
	s.echo.GET("/login", s.handleLoginGet)
	s.echo.POST("/login", s.handleLoginPost)
}

func (s *Server) handleLoginPost(c echo.Context) error {
	var user UserLoginDTO
	err := c.Bind(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	logrus.WithField("user", user).Info("User")
	if err = c.Validate(user); err != nil {
		// we can use this info to build out custom error messages
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	// ok - it worked. so... what now?

	return c.String(http.StatusOK, "Post complete")
}

func (s *Server) handleLoginGet(c echo.Context) error {
	component := views.Login()
	base := views.Base(component)
	templ.Handler(base).ServeHTTP(c.Response().Writer, c.Request())
	return nil
}
