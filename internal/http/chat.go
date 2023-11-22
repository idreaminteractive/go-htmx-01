package http

import (
	"main/internal/views"
	"main/internal/views/dto"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) handleChatByIdPost(c echo.Context) error {

	// creates a new message in the thread
	var message dto.ChatMessageDTO
	if err := c.Bind(&message); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	// var formErrors validation.Errors
	if err := message.Validate(); err != nil {
		// formErrors = err.(validation.Errors)
		// hx return here
		component := views.ChatMessageInput(views.ChatMessageInputProps{PreviousMessage: message.Message, Error: err})
		// render w/ hx
		c.Response().Header().Set("HX-Retarget", "messageInput")
		c.Response().Header().Set("HX-Reswap", "outerHTML")

		renderComponent(component, c)

		return nil
	}
	return nil
}

func (s *Server) handleChatByIdGet(c echo.Context) error {
	// load our data
	props := views.ChatScreenProps{
		ActiveConversations: []views.ConversationItemProps{
			{Handle: "Dave", MessageText: "Hello", Id: 0},
			{Handle: "Dave1", MessageText: "Hello2", Id: 1},
		},
		ActiveChatId: 0,
		CurrentMessages: []views.ChatMessageProps{
			{MessageText: "Hi there", TimeStamp: "10:10 AM"},
		},
	}

	component := views.ChatScreen(props)
	base := views.Base(views.BaseData{Body: component, CSRF: getCSRFValueFromContext(c), Title: "Login"})
	renderComponent(base, c)
	return nil
}

func (s *Server) handleChatGet(c echo.Context) error {
	// this will list our chat maessage

	// ok - this is the root page, so nothing active.

	props := views.ChatScreenProps{
		ActiveConversations: []views.ConversationItemProps{
			{Handle: "Dave", MessageText: "Hello", Id: 0},
			{Handle: "Dave1", MessageText: "Hello2", Id: 1},
		},
		ActiveChatId:    -1,
		CurrentMessages: []views.ChatMessageProps{},
	}

	component := views.ChatScreen(props)
	base := views.Base(views.BaseData{Body: component, CSRF: getCSRFValueFromContext(c), Title: "Login"})
	renderComponent(base, c)
	return nil
}
