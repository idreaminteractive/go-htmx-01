package main

import (
	"context"
	"log"
	"main/components"
	"net/http"

	// "github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	component := components.Layout("Dave")
	component.Render(context.Background(), w)

}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", handleIndex)
	http.ListenAndServe(":3000", r)

}
