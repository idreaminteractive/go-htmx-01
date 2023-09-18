package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	todos "main/db"

	"github.com/go-faker/faker/v4"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	ctx := context.Background()

	db, err := sql.Open("sqlite3", "/litefs/potato.db")
	if err != nil {
		log.Fatal(err)
	}

	queries := todos.New(db)
	a := todos.CreateUserParams{}

	err = faker.FakeData(&a)
	if err != nil {
		fmt.Println(err)
	}
	user, err := queries.CreateUser(ctx, a)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(user)

	// _, err = queries.CreateTodo(ctx, todos.CreateTodoParams{Description: "This is a non-finished", UserID: }
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// finished_todo, err := queries.CreateTodo(ctx, "Finished todo")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = queries.SetTodoDone(ctx, finished_todo.ID)
	// if err != nil {
	// 	log.Fatal(err)
	// }

}
