package app

import (
	"github.com/a-h/templ"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"main/components"
	"main/types"
	"net/http"
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
	// get our data
	getIndex := func() (types.IndexData, error) {
		return types.IndexData{

			Name:  "Index",
			Posts: make([]types.PostData, 0),
		}, nil

	}
	return IndexHandler{
		Log:      log.Default(),
		GetIndex: getIndex,
	}
}
