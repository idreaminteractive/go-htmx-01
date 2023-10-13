package http

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func renderComponent(component templ.Component, c echo.Context) {
	templ.Handler(component).ServeHTTP(c.Response().Writer, c.Request())
}
