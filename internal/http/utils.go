package http

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func renderComponent(component templ.Component, c echo.Context, status ...int) {
	// w/ htmx, we can send back diff statuses, but we'd need custom js to translate them
	// properly into different values
	httpStatus := 200
	if len(status) > 0 {
		httpStatus = status[0]
	}
	c.Response().Writer.WriteHeader(httpStatus)
	templ.Handler(component).ServeHTTP(c.Response().Writer, c.Request())
}
