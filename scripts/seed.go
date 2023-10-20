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
	ctx := context.Background()

	database, err := sql.Open("sqlite3", "/litefs/potato.db")
	if err != nil {
		log.Fatal(err)
	}

	queries := db.New(database)
	as := services.AuthenticationService{Queries: queries}
	user, _ := as.Authenticate(dto.UserLoginDTO{Email: "dwiper@gmail.com", Password: "dave"})

	item, err := queries.CreateNote(ctx, db.CreateNoteParams{UserID: user.ID, Content: faker.Paragraph(), IsPublic: false})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v", item)
	// finished_todo, err := queries.CreateTodo(ctx, "Finished todo")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = queries.SetTodoDone(ctx, finished_todo.ID)
	// if err != nil {
	// 	log.Fatal(err)
	// }

}
