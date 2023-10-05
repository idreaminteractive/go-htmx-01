package main

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"main/config"
	"main/http"
	"main/sqlite"
	"os"
	"os/signal"
)

// See the wtf project for reference
type Program struct {
	DB *sqlite.DB

	HTTPServer *http.Server
	// anything else?

}

func NewProgram() *Program {
	config := config.Parse()
	if config.DopplerConfig != "dev_local" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
	logrus.Info("Logging Ready to Go")
	return &Program{
		// wrapper for our sqlite db functionality
		DB: sqlite.NewDB(config.DatabaseFileName),
		// wrapper for our http server w/ all the services
		HTTPServer: http.NewServer(config),
	}
}

// Close gracefully stops the program.
func (m *Program) Close() error {
	if m.HTTPServer != nil {
		if err := m.HTTPServer.Close(); err != nil {
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
	// todo - fill in the actual run program

	// todo - initialize all the services backed by db

	// todo - attach all services on m.HTTPServer

	if err := m.HTTPServer.Open(":8080"); err != nil {
		return err
	}
	return nil
}

func main() {

	// setup signal handlers

	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() { <-c; cancel() }()

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
