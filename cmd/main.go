package main

import (
	"context"
	"fmt"
	"main/internal/config"
	"main/internal/db"
	"main/internal/http"
	"main/internal/sqlite"
	"os"
	"os/signal"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

// See the wtf project for reference
type Program struct {
	DB     *sqlite.DB
	Port   int
	Server *http.Server
}

func NewProgram() *Program {
	config := config.Parse()
	if config.DopplerConfig != "dev_local" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
	fmt.Printf("%+v", config)

	database := sqlite.NewDB(config.DatabaseFileName)
	if err := database.Open(); err != nil {
		logrus.WithError(err).Error("Of")
		logrus.Panic("Could not open db")
	}

	queries := db.New(database.Connection)

	port, err := strconv.Atoi(config.GoPort)
	if err != nil {
		logrus.Panic("Bad port!")
	}
	return &Program{
		Port:   port,
		DB:     sqlite.NewDB(config.DatabaseFileName),
		Server: http.NewServer(config, queries),
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
