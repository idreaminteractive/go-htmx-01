package http

import (
	"fmt"
	"main/internal/views"
	"main/internal/views/dto"
	"net/http"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (s *Server) handleChatNewPost(c echo.Context) error {

	// make a new chat + redirect to it

	// creates a new message in the thread
	var newChat struct {
		UserId string `form:"userId"`
	}
	if err := c.Bind(&newChat); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	otherUserId, err := strconv.Atoi(newChat.UserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	spew.Dump(newChat)

	//  no errors, message was goot.
	userId := c.Request().Context().Value("userId").(int)
	conv, err := s.services.ChatService.CreateNewConversation(userId, otherUserId)
	spew.Dump(conv)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	c.Response().Header().Set("HX-Redirect", fmt.Sprintf("/chat/%d", conv.ID))
	return nil
}

func (s *Server) handleChatByIdPost(c echo.Context) error {

	chatIdString := c.Param("id")
	chatId, err := strconv.Atoi(chatIdString)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	// creates a new message in the thread
	var message dto.ChatMessageDTO
	if err := c.Bind(&message); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	// var formErrors validation.Errors
	if err := message.Validate(); err != nil {
		formErrors := err.(validation.Errors)
		// hx return here
		component := views.ChatMessageForm(views.ChatMessageFormProps{ActiveChatId: chatId, PreviousMessage: message.Message, Errors: formErrors})
		// render w/ hx
		c.Response().Header().Set("HX-Retarget", "#messageForm")
		c.Response().Header().Set("HX-Reswap", "outerHTML")

		renderComponent(component, c)

		return nil
	}
	//  no errors, message was goot.
	userId := c.Request().Context().Value("userId").(int)

	// add dat message
	_, err = s.services.ChatService.AddMessageToConversation(userId, chatId, message.Message)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	var currentMessages []views.ChatMessageProps
	messages, err := s.services.ChatService.GetConversationsForUser(userId)
	if err != nil {

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())

	}
	for _, conv := range messages {
		if conv.Id == chatId {
			for _, message := range conv.Messages {
				currentMessages = append(currentMessages, views.ChatMessageProps{
					IsOwn:       message.UserId == userId,
					MessageText: message.Content,
					Handle:      message.Handle,
					UserId:      message.UserId,
					TimeStamp:   message.CreatedAt,
				})
			}

		}
	}

	cap := views.ChatActivityProps{
		ActiveChatId:    chatId,
		CurrentMessages: currentMessages,
	}
	// re-render w/ new datas
	// just render chat activity + only use base data
	component := views.ChatActivity(cap)
	// taerget it
	c.Response().Header().Set("HX-Retarget", "#chatActivity")
	c.Response().Header().Set("HX-Reswap", "outerHTML")

	renderComponent(component, c)
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

	userId := c.Request().Context().Value("userId").(int)

	// no matter what, i need my messages for this chat
	var currentMessages []views.ChatMessageProps
	messages, err := s.services.ChatService.GetConversationsForUser(userId)
	if err != nil {

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())

	}
	for _, conv := range messages {
		if conv.Id == chatId {
			for _, message := range conv.Messages {
				currentMessages = append(currentMessages, views.ChatMessageProps{
					IsOwn:       message.UserId == userId,
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
	// what about users we DON'T have a chat with? let's make it a post thing

	possibles, err := s.services.ChatService.GetUsersWithNoConversation(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var possibleData []views.PossibleConversationItemProps
	for _, p := range possibles {
		possibleData = append(possibleData, views.PossibleConversationItemProps{
			Id:     int(p.ID),
			Handle: p.Handle,
		})
	}

	props := views.ChatScreenProps{
		ActiveConversations:   data,
		PossibleConversations: possibleData,
		ActiveChatId:          chatId,
		CurrentMessages:       cap.CurrentMessages,
	}

	component := views.ChatScreen(props)
	base := views.Base(views.BaseData{Body: component, CSRF: csrfFromRequest(c.Request()), Title: "Login"})
	renderComponent(base, c)
	return nil
}

func (s *Server) getConversationData(c echo.Context) ([]views.ConversationItemProps, error) {

	userId := c.Request().Context().Value("userId").(int)

	// ok - this is the root page, so nothing active.

	data, err := s.services.ChatService.GetConversationsForUser(userId)
	if err != nil {

		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())

	}
	ActiveConversations := []views.ConversationItemProps{}
	for _, conversation := range data {
		// this is bad!
		otherUser, err := s.services.ChatService.GetOtherUserInConversation(userId, conversation.Id)
		if err != nil {
			logrus.Error(err)
			continue
		}
		firstMessage := conversation.Messages[0]
		mText := "No message"
		if firstMessage.UserId == userId {
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
	userId := c.Request().Context().Value("userId").(int)
	possibles, err := s.services.ChatService.GetUsersWithNoConversation(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var possibleData []views.PossibleConversationItemProps
	for _, p := range possibles {
		possibleData = append(possibleData, views.PossibleConversationItemProps{
			Id:     int(p.ID),
			Handle: p.Handle,
		})
	}
	props := views.ChatScreenProps{
		PossibleConversations: possibleData,
		ActiveConversations:   data,
		ActiveChatId:          -1,
		CurrentMessages:       []views.ChatMessageProps{},
	}

	component := views.ChatScreen(props)
	base := views.Base(views.BaseData{Body: component, CSRF: csrfFromRequest(c.Request()), Title: "Login"})
	renderComponent(base, c)
	return nil
}
