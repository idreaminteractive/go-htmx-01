package main

import (
	"context"
	"database/sql"
	"log"
	todos "main/db"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	ctx := context.Background()

	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}

	queries := todos.New(db)

	_, err = queries.CreateTodo(ctx, "Unfinished todo")
	if err != nil {
		log.Fatal(err)
	}

	finished_todo, err := queries.CreateTodo(ctx, "Finished todo")
	if err != nil {
		log.Fatal(err)
	}

	err = queries.SetTodoDone(ctx, finished_todo.ID)
	if err != nil {
		log.Fatal(err)
	}

}
