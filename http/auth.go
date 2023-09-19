package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type UserLoginDTO struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

func (s *Server) registerAuthRoutes() {
	s.echo.GET("/login", s.handleLoginGet)
	s.echo.POST("/login", s.handleLoginPost)
}

func (s *Server) handleLoginPost(c echo.Context) error {
	var user UserLoginDTO
	err := c.Bind(&user)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	logrus.WithField("user", user).Info("User")
	return c.String(http.StatusOK, "polsted")
}

func (s *Server) handleLoginGet(c echo.Context) error {

	return c.String(http.StatusOK, "login")

}
