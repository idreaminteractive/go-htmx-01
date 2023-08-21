//go:generate go get -u github.com/valyala/quicktemplate/qtc
//go:generate qtc -dir=views

package main

import (
	"html/template"
	"log"
	todos "main/db"
	"main/views"
	"net/http"

	"context"
	"database/sql"
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	fmt.Printf("%s\n", views.Hello("Foo"))
	fmt.Printf("%s\n", views.Hello("potato"))

	ctx := context.Background()

	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}

	queries := todos.New(db)

	todos, err := queries.ListTodos(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(todos)

	r := chi.NewRouter()

	const tpl = `
	<!DOCTYPE html>
	<html>
		<head>
		<script src="https://cdn.jsdelivr.net/npm/@unocss/runtime/uno.global.js"></script>
			<meta charset="UTF-8">
			<title>{{.Title}}</title>
		</head>
		<body>
		<div class="h-full text-center flex select-none all:transition-400"> Potatol </div>
		<div class="text-blue-500">blue</div>
			{{range .Items}}<div>{{ . }}</div>{{else}}<div><strong>no rows</strong></div>{{end}}
		</body>
	</html>`
	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}
	t, err := template.New("webpage").Parse(tpl)

	// use t?

	data := struct {
		Title string
		Items []string
	}{
		Title: "My page",
		Items: []string{
			"My photos",
			"My blog",
		},
	}

	check(err)
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		t.Execute(w, data)
	})
	http.ListenAndServe(":3000", r)

}
