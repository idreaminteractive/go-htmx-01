package main

import (
	"context"
	"fmt"
	"log/slog"
	"main/internal/config"
	"main/internal/db"
	"main/internal/http"
	"main/internal/sqlite"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/go-chi/httplog/v2"
	_ "github.com/mattn/go-sqlite3"
)

// See the wtf project for reference
type Program struct {
	DB     *sqlite.DB
	Port   int
	Server *http.Server
}

func NewProgram() *Program {
	config := config.Parse()

	// the logger will be available
	// via both the controller on the context
	// as well as in the server struct
	logger := httplog.NewLogger("go-htmx-01", httplog.Options{
		// JSON:             true,
		LogLevel:         slog.LevelDebug,
		Concise:          true,
		RequestHeaders:   true,
		MessageFieldName: "message",
		// TimeFieldFormat: time.RFC850,
		Tags: map[string]string{
			"version": "v1.0-81aa4244d9fc8076a",
			"env":     config.DopplerConfig,
		},
		QuietDownRoutes: []string{
			"/healthz",
		},
		QuietDownPeriod: 10 * time.Second,
		// SourceFieldName: "source",
	})

	database := sqlite.NewDB(config.DatabaseFileName)
	if err := database.Open(); err != nil {
		logger.Error("Could not open db", err)

		panic("Could not open db")
	}

	queries := db.New(database.Connection)

	port, err := strconv.Atoi(config.GoPort)
	if err != nil {
		panic("Bad port!")
	}
	return &Program{
		Port:   port,
		DB:     sqlite.NewDB(config.DatabaseFileName),
		Server: http.NewServer(config, queries, logger),
	}
}

// Close gracefully stops the program.
func (m *Program) Close() error {
	if m.Server != nil {
		if err := m.Server.Close(); err != nil {
			return err
		}
	}
	if m.DB != nil {
		if err := m.DB.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (m *Program) Run(ctx context.Context) error {
	// passing port does nothing here.
	if err := m.Server.Open(fmt.Sprintf(":%d", m.Port)); err != nil {
		return err
	}
	return nil
}

func main() {

	// setup signal handlers

	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	// if we receive an os interrupt, we run cancel which kills the rest of it all
	signal.Notify(c, os.Interrupt)
	go func() { <-c; cancel() }()

	// sets it all up
	m := NewProgram()

	// Execute program.
	if err := m.Run(ctx); err != nil {
		m.Close()
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Wait for CTRL-C.
	<-ctx.Done()

	// Clean up program.
	if err := m.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
