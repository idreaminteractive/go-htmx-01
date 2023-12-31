package http

import (
	"fmt"

	base "main/internal/views/base"
	chat "main/internal/views/chat"
	flash "main/internal/views/components/flash"
	"main/internal/views/dto"
	"net/http"
	"strconv"

	"github.com/angelofallars/htmx-go"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
)

type newChat struct {
	UserId string `form:"userId"`
}

func (nc newChat) Bind(r *http.Request) error {

	return nil
}
func (s *Server) handleChatNewPost(w http.ResponseWriter, r *http.Request) {

	// make a new chat + redirect to it

	// creates a new message in the thread
	var nc newChat

	if err := render.Bind(r, &nc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	otherUserId, err := strconv.Atoi(nc.UserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//  no errors, message was goot.
	userId := s.getUserIdFromCTX(r)
	conv, err := s.services.ChatService.CreateNewConversation(userId, otherUserId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	htmx.NewResponse().Redirect(fmt.Sprintf("/chat/%d", conv.ID)).Write(w)

}

func (s *Server) handleChatByIdPost(w http.ResponseWriter, r *http.Request) {

	chatIdString := chi.URLParam(r, "id")
	chatId, err := strconv.Atoi(chatIdString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// creates a new message in the thread
	var message dto.ChatMessageDTO
	if err := render.Bind(r, &message); err != nil {
		// general render bind error
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// var formErrors validation.Errors
	if err := message.Validate(); err != nil {
		formErrors := err.(validation.Errors)
		// hx return here
		component := chat.ChatMessageForm(chat.ChatMessageFormProps{ActiveChatId: chatId, PreviousMessage: message.Message, Errors: formErrors})
		// render w/ hx
		htmx.NewResponse().Retarget("#messageForm").Reswap(htmx.SwapOuterHTML).RenderTempl(r.Context(), w, component)

		return
	}
	//  no errors, message was goot.
	userId := s.getUserIdFromCTX(r)

	// add dat message
	_, err = s.services.ChatService.AddMessageToConversation(userId, chatId, message.Message)
	if err != nil {
		fmt.Println("We are in error oof. why no return?")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// do sse for it
	count, err := s.services.ChatService.GetTotalMessagCount()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.services.SSEEventBus.Send("message-count", fmt.Sprintf("%d", count))

	var currentMessages []chat.ChatMessageProps
	messages, err := s.services.ChatService.GetConversationsForUser(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

	for _, conv := range messages {
		if conv.Id == chatId {
			for _, message := range conv.Messages {
				currentMessages = append(currentMessages, chat.ChatMessageProps{
					IsOwn:       message.UserId == userId,
					MessageText: message.Content,
					Handle:      message.Handle,
					UserId:      message.UserId,
					TimeStamp:   message.CreatedAt,
				})
			}

		}
	}

	cap := chat.ChatActivityProps{
		ActiveChatId:    chatId,
		CurrentMessages: currentMessages,
	}
	// re-render w/ new datas
	// just render chat activity + only use base data
	component := chat.ChatActivity(cap)

	// for this, it's not a real flash message,
	// but an oob swap into the page to simulate it.

	htmx.NewResponse().
		Retarget("#chatActivity").
		Reswap(htmx.SwapOuterHTML).
		RenderTempl(r.Context(), w, flash.FlashWrapper(component, flash.FlashProps{
			Message: "Hi there!",
		}))

}

func (s *Server) handleChatByIdGet(w http.ResponseWriter, r *http.Request) {
	// load our data
	// this could be a lot leaner

	chatIdString := chi.URLParam(r, "id")
	chatId, err := strconv.Atoi(chatIdString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	userId := s.getUserIdFromCTX(r)

	// no matter what, i need my messages for this chat
	var currentMessages []chat.ChatMessageProps
	messages, err := s.services.ChatService.GetConversationsForUser(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}
	for _, conv := range messages {
		if conv.Id == chatId {
			for _, message := range conv.Messages {
				currentMessages = append(currentMessages, chat.ChatMessageProps{
					IsOwn:       message.UserId == userId,
					MessageText: message.Content,
					Handle:      message.Handle,
					UserId:      message.UserId,
					TimeStamp:   message.CreatedAt,
				})
			}

		}
	}

	cap := chat.ChatActivityProps{
		ActiveChatId:    chatId,
		CurrentMessages: currentMessages,
	}
	if htmx.IsHTMX(r) {
		// just render chat activity + only use base data
		component := chat.ChatActivity(cap)
		htmx.NewResponse().RenderTempl(r.Context(), w, component)
		return
	}
	data, err := s.getConversationData(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// what about users we DON'T have a chat with? let's make it a post thing

	possibles, err := s.services.ChatService.GetUsersWithNoConversation(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var possibleData []chat.PossibleConversationItemProps
	for _, p := range possibles {
		possibleData = append(possibleData, chat.PossibleConversationItemProps{
			Id:     int(p.ID),
			Handle: p.Handle,
		})
	}

	props := chat.ChatScreenProps{
		ActiveConversations:   data,
		PossibleConversations: possibleData,
		ActiveChatId:          chatId,
		CurrentMessages:       cap.CurrentMessages,
	}

	component := chat.ChatScreen(props)
	base := base.Base(base.BaseData{Body: component, CSRF: csrfFromRequest(r), Title: "Chatting"})
	htmx.NewResponse().RenderTempl(r.Context(), w, base)

}

func (s *Server) getConversationData(userId int) ([]chat.ConversationItemProps, error) {

	// ok - this is the root page, so nothing active.

	data, err := s.services.ChatService.GetConversationsForUser(userId)
	if err != nil {

		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())

	}
	ActiveConversations := []chat.ConversationItemProps{}
	for _, conversation := range data {
		// this is bad!
		otherUser, err := s.services.ChatService.GetOtherUserInConversation(userId, conversation.Id)
		if err != nil {
			s.logger.Error("Error getting other user", err)
			continue
		}
		firstMessage := conversation.Messages[0]
		mText := "No message"
		if firstMessage.UserId == userId {
			mText = "> " + firstMessage.Content
		} else {
			mText = "< " + firstMessage.Content
		}
		ActiveConversations = append(ActiveConversations, chat.ConversationItemProps{
			Id: conversation.Id, Handle: otherUser.Handle, MessageText: mText,
		})
	}
	return ActiveConversations, nil
}

func (s *Server) handleChatGet(w http.ResponseWriter, r *http.Request) {
	// this will list our chat maessage

	// get our convos +
	userId := s.getUserIdFromCTX(r)
	data, err := s.getConversationData(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	possibles, err := s.services.ChatService.GetUsersWithNoConversation(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var possibleData []chat.PossibleConversationItemProps
	for _, p := range possibles {
		possibleData = append(possibleData, chat.PossibleConversationItemProps{
			Id:     int(p.ID),
			Handle: p.Handle,
		})
	}
	props := chat.ChatScreenProps{
		PossibleConversations: possibleData,
		ActiveConversations:   data,
		ActiveChatId:          -1,
		CurrentMessages:       []chat.ChatMessageProps{},
	}

	component := chat.ChatScreen(props)
	base := base.Base(base.BaseData{Body: component, CSRF: csrfFromRequest(r), Title: "Login"})
	htmx.NewResponse().RenderTempl(r.Context(), w, base)

}
