package main

import (
	"context"
	"fmt"
	"main/http"
	"main/sqlite"
	"os"
	"os/signal"

	"github.com/caarlos0/env/v9"
)

type EnvConfig struct {
	DatabaseFileName string `env:"DATABASE_FILENAME" envDefault:"/litefs/potato.db"`
	GoPort           string `env:"GO_PORT" envDefault:"8080"`
}

// See the wtf project for reference
type Program struct {
	Config EnvConfig
	DB     *sqlite.DB

	HTTPServer *http.Server
	// anything else?

}

func NewProgram() *Program {
	config := EnvConfig{}
	err := env.Parse(&config)
	if err != nil {
		panic("Could not parse env")
	}

	return &Program{
		Config: config,
		// wrapper for our sqlite db functionality
		DB: sqlite.NewDB(config.DatabaseFileName),
		// wrapper for our http server w/ all the services
		HTTPServer: http.NewServer(),
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
	fmt.Println("Listening on " + m.Config.GoPort)
	if err := m.HTTPServer.Open(":" + m.Config.GoPort); err != nil {
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
