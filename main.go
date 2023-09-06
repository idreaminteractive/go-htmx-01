package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	todos "main/db"
	"strconv"

	"github.com/caarlos0/env/v9"
	_ "github.com/mattn/go-sqlite3"

	"net/http"
	"time"

	g "github.com/maragudk/gomponents"
	hx "github.com/maragudk/gomponents-htmx"
	hxhttp "github.com/maragudk/gomponents-htmx/http"
	c "github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/html"
	ghttp "github.com/maragudk/gomponents/http"
)

var myTodos *todos.Queries

type EnvConfig struct {
	DatabaseFileName string `env:"DATABASE_FILENAME" envDefault:"/litefs/potato.db"`
	GoPort           string `env:"GO_PORT" envDefault:"8080"`
}

func main() {

	config := EnvConfig{}
	if err := env.Parse(&config); err != nil {
		fmt.Printf("%+v\n", err)
	}

	db, err := sql.Open("sqlite3", config.DatabaseFileName)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(`PRAGMA journal_mode = wal;`); err != nil {
		log.Fatalf("enable wal: %w", err)
	}

	// Enable foreign key checks. For historical reasons, SQLite does not check
	// foreign key constraints by default... which is kinda insane. There's some
	// overhead on inserts to verify foreign key integrity but it's definitely
	// worth it.
	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		log.Fatalf("foreign keys pragma: %w", err)
	}
	//

	myTodos = todos.New(db)

	if err := start(config.GoPort); err != nil {
		log.Fatalln("Error:", err)
	}

}

func start(port string) error {
	now := time.Now()

	mux := http.NewServeMux()
	mux.HandleFunc("/", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (g.Node, error) {
		if r.Method == http.MethodPost && hxhttp.IsBoosted(r.Header) {
			now = time.Now()

			hxhttp.SetPushURL(w.Header(), "/?time="+now.Format(timeFormat))

			return partial(now), nil
		}
		return page(now), nil
	}))

	log.Println("Starting on Port " + port)
	if err := http.ListenAndServe(":"+port, mux); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

const timeFormat = "15:04:05"

func page(now time.Time) g.Node {
	ctx := context.Background()
	todoList, err := myTodos.ListTodos(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return c.HTML5(c.HTML5Props{
		Title: now.Format(timeFormat),
		Head: []g.Node{
			Script(Src("https://cdn.tailwindcss.com?plugins=forms,typography")),
			Script(Src("https://unpkg.com/htmx.org")),
		},
		Body: []g.Node{
			Div(Class("max-w-7xl mx-auto p-4 prose lg:prose-lg xl:prose-xl"),
				Div(Class("text-lg"), g.Text("I am a new div")),
				H1(g.Text(`gomponents + HTMX`)),
				P(g.Textf(`Time at last full page refresh was %v.`, now.Format(timeFormat))),
				partial(now),
				FormEl(Method("post"), Action("/"), hx.Boost("true"), hx.Target("#partial"), hx.Swap("outerHTML"),
					Button(Type("submit"), g.Text(`Update time`),
						Class("rounded-md border border-transparent bg-orange-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-orange-700 focus:outline-none focus:ring-2 focus:ring-orange-500 focus:ring-offset-2"),
					),
				),
				Ul(
					Div(
						g.Group(
							g.Map(todoList, func(todo todos.Todo) g.Node {
								return Li(g.Text(strconv.FormatInt(todo.ID, 10) + ":" + todo.Description))
							}),
						),
					),
				),
			),
		},
	})
}

func partial(now time.Time) g.Node {
	return P(ID("partial"), g.Textf(`Time was last updated at %v.`, now.Format(timeFormat)))
}
