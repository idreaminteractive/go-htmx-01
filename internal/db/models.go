// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package db

import (
	"time"
)

type Conversation struct {
	ID        int64
	CreatedAt time.Time
}

type Message struct {
	ID             int64
	ConversationID int64
	UserID         int64
	Content        string
	CreatedAt      time.Time
}

type Note struct {
	ID        int64
	Content   string
	UserID    int64
	IsPublic  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	ID        int64
	Email     string
	Handle    string
	Password  string
	CreatedAt time.Time
}

type UserConversation struct {
	UserID         int64
	ConversationID int64
}
