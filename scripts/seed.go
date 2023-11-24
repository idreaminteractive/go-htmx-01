package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"main/internal/db"
	"main/internal/services"
	"main/internal/views/dto"

	"github.com/go-faker/faker/v4"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	database, err := sql.Open("sqlite3", "/litefs/potato.db")
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	queries := db.New(database)

	as := services.InitAuthService(&services.ServiceLocator{}, queries)

	one, err := as.Register(dto.RegisterDTO{Handle: "Dave", Email: "dwiper@gmail.com", Password: "dave", ConfirmPassword: "dave"})
	if err != nil {
		log.Fatal(err)
	}
	// ok - les register three others + start some chats

	two, err := as.Register(dto.RegisterDTO{Handle: faker.Username(), Email: faker.Email(), Password: "dave", ConfirmPassword: "dave"})
	if err != nil {
		log.Fatal(err)
	}
	three, err := as.Register(dto.RegisterDTO{Handle: faker.Username(), Email: faker.Email(), Password: "dave", ConfirmPassword: "dave"})
	if err != nil {
		log.Fatal(err)
	}
	_, err = as.Register(dto.RegisterDTO{Handle: faker.Username(), Email: faker.Email(), Password: "dave", ConfirmPassword: "dave"})
	if err != nil {
		log.Fatal(err)
	}
	// associate one + two to conversation
	conv, err := queries.CreateConversation(ctx)
	if err != nil {
		log.Fatal(err)
	}
	_, err = queries.LinkUserToConversation(ctx, db.LinkUserToConversationParams{UserID: one.ID, ConversationID: conv.ID})
	if err != nil {
		log.Fatal(err)
	}
	_, err = queries.LinkUserToConversation(ctx, db.LinkUserToConversationParams{UserID: two.ID, ConversationID: conv.ID})
	if err != nil {
		log.Fatal(err)
	}
	// add a message or two
	_, err = queries.CreateMessage(ctx, db.CreateMessageParams{UserID: one.ID, ConversationID: conv.ID, Content: "Message from Dave"})
	if err != nil {
		log.Fatal(err)
	}
	_, err = queries.CreateMessage(ctx, db.CreateMessageParams{UserID: two.ID, ConversationID: conv.ID, Content: "Message TO Dave"})
	if err != nil {
		log.Fatal(err)
	}

	// create second convo w/ dave + threee
	conv, err = queries.CreateConversation(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nSecond one %d\n", conv.ID)
	_, err = queries.LinkUserToConversation(ctx, db.LinkUserToConversationParams{UserID: one.ID, ConversationID: conv.ID})
	if err != nil {
		log.Fatal(err)
	}
	_, err = queries.LinkUserToConversation(ctx, db.LinkUserToConversationParams{UserID: three.ID, ConversationID: conv.ID})
	if err != nil {
		log.Fatal(err)
	}
	// add a message or two
	_, err = queries.CreateMessage(ctx, db.CreateMessageParams{UserID: one.ID, ConversationID: conv.ID, Content: "Hi Three, I'm dave!"})
	if err != nil {
		log.Fatal(err)
	}
	_, err = queries.CreateMessage(ctx, db.CreateMessageParams{UserID: three.ID, ConversationID: conv.ID, Content: "Message from Three to Dave"})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n\nDone seeding!")
}
