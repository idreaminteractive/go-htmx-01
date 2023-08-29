package main

import (
	"log"
	"main/app"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)

	r.Get("/", app.NewIndexHandler().ServeHTTP)

	http.ListenAndServe(":3000", r)

}
