package app

import (
	"log"
	"main/components"
	"main/types"
	"net/http"

	"github.com/a-h/templ"
	_ "github.com/mattn/go-sqlite3"
)

type IndexHandler struct {
	Log      *log.Logger
	GetIndex func() (types.IndexData, error)
}

func (ih IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ps, err := ih.GetIndex()
	if err != nil {
		ih.Log.Printf("failed to get index: %v", err)
		http.Error(w, "failed to retrieve index", http.StatusInternalServerError)
		return
	}
	templ.Handler(components.IndexTemplate(ps)).ServeHTTP(w, r)
}

// Actual data object for index of the site
func NewIndexHandler() IndexHandler {
	// get our data from our database here.
	getIndex := func() (types.IndexData, error) {
		return types.IndexData{
			Name: "Index",
			Posts: []types.PostData{
				{Id: 1, Name: "Potatos", Content: "This is the post content"},
				{Id: 2, Name: "Tomatoes", Content: "This is the post contentß"},
			},
		}, nil

	}
	return IndexHandler{
		Log:      log.Default(),
		GetIndex: getIndex,
	}
}
