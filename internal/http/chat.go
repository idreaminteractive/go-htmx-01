package http

import (
	"main/internal/views"

	"github.com/labstack/echo/v4"
)

func (s *Server) handleChatGet(c echo.Context) error {
	// this will list our chat maessage


	
	component := views.ChatScreen()
	base := views.Base(views.BaseData{Body: component, CSRF: getCSRFValueFromContext(c), Title: "Login"})
	renderComponent(base, c)
	return nil
}
