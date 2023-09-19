package http

import (
	"main/views/routes"
	"net/http"

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
	return c.String(http.StatusOK, "Post complete")
}

func (s *Server) handleLoginGet(c echo.Context) error {
	return routes.LoginPage().Render(c.Response().Writer)

}
