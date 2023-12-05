package services

import (
	"context"
	"encoding/json"

	"main/internal/db"

	"github.com/go-chi/httplog/v2"
)

type ChatService struct {
	sl             *ServiceLocator
	queries        *db.Queries
	MessageChannel chan []byte
	logger         *httplog.Logger
}

func InitChatService(sl *ServiceLocator, queries *db.Queries, logger *httplog.Logger) *ChatService {
	return &ChatService{
		sl:             sl,
		queries:        queries,
		logger:         logger,
		MessageChannel: make(chan []byte),
	}
}

type ConversationMessages struct {
	MessageId int    `json:"message_id"`
	Content   string `json:"content"`
	UserId    int    `json:"user_id"`
	Handle    string `json:"handle"`
	CreatedAt string `json:"created_at"`
}

type Conversation struct {
	Id       int
	Handle   string
	UserId   int
	Messages []ConversationMessages
}

func (cs *ChatService) GetConversationsForUser(userId int) ([]Conversation, error) {
	ctx := context.Background()
	conversations, err := cs.queries.GetConversationsList(ctx, int64(userId))
	if err != nil {
		cs.logger.Error("Error gettting conversations", err)
		// user exists already
		return nil, &Error{Code: EINTERNAL, Message: err.Error()}
	}

	// ok - so our return will actually be conversation ID + then the return in a struct
	// note, some stuff may be fasyter, but this is likely the safest
	output := []Conversation{}
	for _, c := range conversations {
		var cm []ConversationMessages
		err := json.Unmarshal([]byte(c.ConversationMessages.(string)), &cm)
		if err != nil {
			cs.logger.Error("Err unmarshaling", err)
			continue
		}

		output = append(output, Conversation{
			Id:       int(c.ConversationID),
			Handle:   c.Handle,
			UserId:   int(c.UserID),
			Messages: cm,
		})

	}

	return output, nil

}

type otherUserReturn struct {
	Handle string
	Id     int
}

func (cs *ChatService) CreateNewConversation(userIdOne, userIdTwo int) (*db.Conversation, error) {
	ctx := context.Background()
	// how do i get a txn here?
	conv, err := cs.queries.CreateConversation(ctx)
	if err != nil {
		cs.logger.Error("Error creating conversation", err)
		// user exists already
		return nil, &Error{Code: EINTERNAL, Message: err.Error()}
	}

	_, err = cs.queries.LinkUserToConversation(ctx, db.LinkUserToConversationParams{UserID: int64(userIdOne), ConversationID: conv.ID})
	if err != nil {
		cs.logger.Error("Error linking", err)

		return nil, &Error{Code: EINTERNAL, Message: err.Error()}
	}
	_, err = cs.queries.LinkUserToConversation(ctx, db.LinkUserToConversationParams{UserID: int64(userIdTwo), ConversationID: conv.ID})
	if err != nil {
		cs.logger.Error("Error linking", err)

		return nil, &Error{Code: EINTERNAL, Message: err.Error()}
	}

	// also add in a message for good measure!
	_, err = cs.queries.CreateMessage(ctx, db.CreateMessageParams{UserID: int64(userIdOne), ConversationID: conv.ID, Content: "Hello there!"})
	if err != nil {
		cs.logger.Error("Error creating message", err)

		return nil, &Error{Code: EINTERNAL, Message: err.Error()}
	}

	return &conv, nil
}

func (cs *ChatService) GetOtherUserInConversation(userId, conversationId int) (*otherUserReturn, error) {
	ctx := context.Background()
	data, err := cs.queries.GetOtherConversationUser(ctx, db.GetOtherConversationUserParams{ConversationID: int64(conversationId), ID: int64(userId)})
	if err != nil {
		cs.logger.Error("Error getting other user", err)
		return nil, &Error{Code: EINTERNAL, Message: err.Error()}
	}
	return &otherUserReturn{Handle: data.Handle, Id: int(data.ID)}, nil
}

func (cs *ChatService) AddMessageToConversation(userId, conversationId int, content string) (*db.Message, error) {
	ctx := context.Background()

	msg, err := cs.queries.CreateMessage(ctx, db.CreateMessageParams{UserID: int64(userId), ConversationID: int64(conversationId), Content: content})
	if err != nil {
		cs.logger.Error("Error creating message", err)
		return nil, &Error{Code: EINTERNAL, Message: err.Error()}
	}

	return &msg, nil
}

func (cs *ChatService) GetTotalMessagCount() (int64, error) {
	ctx := context.Background()

	count, err := cs.queries.GetTotalNumMessages(ctx)
	if err != nil {
		cs.logger.Error("Error getting number of messages", err)
		return 0, &Error{Code: EINTERNAL, Message: err.Error()}
	}
	// we made dat message
	return count, nil
}

func (cs *ChatService) GetUsersWithNoConversation(userId int) ([]db.PossibleConversationUsersRow, error) {
	ctx := context.Background()

	// it's all the same args, i feel like there's a better way... lol
	possibles, err := cs.queries.PossibleConversationUsers(ctx, int64(userId))
	if err != nil {
		cs.logger.Error("Error getting possible user conversions", err)
		return nil, &Error{Code: EINTERNAL, Message: err.Error()}
	}

	return possibles, nil
}
