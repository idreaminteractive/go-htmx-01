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
	r.Use(middleware.RequestID)

	r.Get("/", handleIndex)

	http.ListenAndServe(":3000", r)

}

type IndexHandler struct {
	Log      *log.Logger
	GetIndex func() (IndexData, error)
}

func (ih IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ps, err := ih.GetIndex()
	if err != nil {
		ih.Log.Printf("failed to get index: %v", err)
		http.Error(w, "failed to retrieve index", http.StatusInternalServerError)
		return
	}
	templ.Handler(IndexTemplate(ps)).ServeHTTP(w, r)
}

// Actual data object for index of the site
type IndexData struct {
	name  string
	posts []PostData
}

func NewIndexHandler() IndexHandler {
	// get our data
	getIndex := func() (IndexData, error) {
		return IndexData{
			name:  "Index",
			posts: make([]PostData, 0),
		}, nil

	}
	return IndexHandler{
		Log:      log.Default(),
		GetIndex: getIndex,
	}
}

type PostHandler struct {
	Log     *log.Logger
	GetPost func(postId int) (PostData, error)
}

type PostData struct {
	id      int
	name    string
	content string
}
