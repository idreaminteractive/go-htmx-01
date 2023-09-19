package http

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
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
		// is it better to do this in a service, or at controller?
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println()
			fmt.Println(err.Namespace())
			fmt.Println(err.Field())
			fmt.Println(err.StructNamespace())
			fmt.Println(err.StructField())
			fmt.Println(err.Tag())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println()
		}
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, "polsted")
}

func (s *Server) handleLoginGet(c echo.Context) error {

	return c.String(http.StatusOK, "login")

}
