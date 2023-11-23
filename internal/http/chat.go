package http

import (
	"main/internal/views"
	"main/internal/views/dto"
	"net/http"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (s *Server) handleChatByIdPost(c echo.Context) error {

	// creates a new message in the thread
	var message dto.ChatMessageDTO
	if err := c.Bind(&message); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	// var formErrors validation.Errors
	if err := message.Validate(); err != nil {
		formErrors := err.(validation.Errors)
		// hx return here
		component := views.ChatMessageForm(views.ChatMessageFormProps{PreviousMessage: message.Message, Errors: formErrors})
		// render w/ hx
		c.Response().Header().Set("HX-Retarget", "#messageForm")
		c.Response().Header().Set("HX-Reswap", "outerHTML")

		renderComponent(component, c)

		return nil
	}
	return nil
}

func (s *Server) handleChatByIdGet(c echo.Context) error {
	// load our data
	// this could be a lot leaner

	chatIdString := c.Param("id")
	chatId, err := strconv.Atoi(chatIdString)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	userId := c.Request().Context().Value("userId")

	// no matter what, i need my messages for this chat
	var currentMessages []views.ChatMessageProps
	messages, err := s.services.ChatService.GetConversationsForUser(userId.(int))
	if err != nil {

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())

	}
	for _, conv := range messages {
		if conv.Id == chatId {
			for _, message := range conv.Messages {
				currentMessages = append(currentMessages, views.ChatMessageProps{
					MessageText: message.Content,
					Handle:      message.Handle,
					UserId:      message.UserId,
					TimeStamp:   message.CreatedAt,
				})
			}

		}
	}

	isHX := c.Request().Header.Get("Hx-Request")
	cap := views.ChatActivityProps{
		ActiveChatId:    chatId,
		CurrentMessages: currentMessages,
	}
	if isHX == "true" {
		// just render chat activity + only use base data
		component := views.ChatActivity(cap)
		renderComponent(component, c)
		return nil
	}
	data, err := s.getConversationData(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	props := views.ChatScreenProps{
		ActiveConversations: data,
		ActiveChatId:        chatId,
		CurrentMessages:     cap.CurrentMessages,
	}

	component := views.ChatScreen(props)
	base := views.Base(views.BaseData{Body: component, CSRF: getCSRFValueFromContext(c), Title: "Login"})
	renderComponent(base, c)
	return nil
}

func (s *Server) getConversationData(c echo.Context) ([]views.ConversationItemProps, error) {
	sess, err := s.services.SessionService.ReadSession(c)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Could not read session")

	}
	// ok - this is the root page, so nothing active.

	data, err := s.services.ChatService.GetConversationsForUser(sess.UserId)
	if err != nil {

		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())

	}

	ActiveConversations := []views.ConversationItemProps{}
	for _, conversation := range data {
		// this is bad!
		otherUser, err := s.services.ChatService.GetOtherUserInConversation(sess.UserId, conversation.Id)
		if err != nil {
			logrus.Error(err)
			continue
		}
		firstMessage := conversation.Messages[0]
		mText := "No message"
		if firstMessage.UserId == sess.UserId {
			mText = "> " + firstMessage.Content
		} else {
			mText = "< " + firstMessage.Content
		}
		ActiveConversations = append(ActiveConversations, views.ConversationItemProps{
			Id: conversation.Id, Handle: otherUser.Handle, MessageText: mText,
		})
	}
	return ActiveConversations, nil
}

func (s *Server) handleChatGet(c echo.Context) error {
	// this will list our chat maessage
	// get our convos +
	data, err := s.getConversationData(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	props := views.ChatScreenProps{
		ActiveConversations: data,
		ActiveChatId:        -1,
		CurrentMessages:     []views.ChatMessageProps{},
	}

	component := views.ChatScreen(props)
	base := views.Base(views.BaseData{Body: component, CSRF: getCSRFValueFromContext(c), Title: "Login"})
	renderComponent(base, c)
	return nil
}
