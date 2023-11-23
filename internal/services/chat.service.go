package services

import (
	"context"
	"encoding/json"
	"main/internal/db"

	"github.com/sirupsen/logrus"
)

type ChatService struct {
	sl      *ServiceLocator
	queries *db.Queries
}

func InitChatService(sl *ServiceLocator, queries *db.Queries) *ChatService {
	return &ChatService{
		sl:      sl,
		queries: queries,
	}
}

func (cs *ChatService) StartConversation() {

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
		logrus.Error(err)
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
			logrus.Error(err)
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
