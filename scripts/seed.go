package main

import (
	"database/sql"
	"log"

	"main/internal/db"
	"main/internal/services"
	"main/internal/views/dto"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	database, err := sql.Open("sqlite3", "/litefs/potato.db")
	if err != nil {
		log.Fatal(err)
	}

	queries := db.New(database)

	as := services.InitAuthService(&services.ServiceLocator{}, queries)

	_, err = as.Register(dto.RegisterDTO{Handle: "Dave", Email: "dwiper@gmail.com", Password: "dave", ConfirmPassword: "dave"})

	if err != nil {
		log.Fatal(err)
	}

}
