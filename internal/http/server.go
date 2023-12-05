package http

import (
	"context"
	"encoding/gob"
	"fmt"

	"main/internal/config"
	"main/internal/db"
	"main/internal/hotreload"
	"main/internal/services"
	"main/internal/session"
	"main/internal/sse"

	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/gorilla/csrf"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
	"github.com/gorilla/sessions"
)

const ShutdownTimeout = 1 * time.Second

// i want to
type Server struct {
	router   *chi.Mux
	server   *http.Server
	config   *config.EnvConfig
	services *services.ServiceLocator
	logger   *httplog.Logger
}

type ServerSetupStruct struct {
	SessionSecret string
	// default bool is false, so we generally want to enable it
	DisableCSRF bool
	Logger      *httplog.Logger
	Env         string
}

func setupServer(config ServerSetupStruct) *chi.Mux {
	// sets up echo with standard things
	// we attach it here in order to allow tests to use it as well.
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	// add our logger to chi + to server struct
	r.Use(httplog.RequestLogger(config.Logger))

	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.Compress(5))
	gob.Register(services.SessionPayload{})

	r.Use(session.Middleware(sessions.NewCookieStore([]byte(config.SessionSecret))))
	r.Use(AddEnvMiddleware(config.Env))
	validation.ErrorTag = "form"
	if !config.DisableCSRF {
		csrfMiddleware :=
			csrf.Protect([]byte("32-byte-long-auth-key"))
		r.Use(csrfMiddleware)

	}

	if config.Env == "dev_local" {
		handler := hotreload.New()
		r.Get("/hmr", handler.ServeHTTP)
	}

	return r
}

func NewServer(config *config.EnvConfig, queries *db.Queries, logger *httplog.Logger) *Server {
	// This is where we initialize all our services and attach to our
	// server
	r := setupServer(ServerSetupStruct{Logger: logger, SessionSecret: config.SessionSecret, Env: config.DopplerConfig})

	// create our general production events sse handler
	sseHandler := sse.New()

	// setup our service locator
	sl := services.ServiceLocator{
		SSEEventBus: sseHandler,
	}

	sl.AuthenticationService = services.InitAuthService(&sl, queries, logger)
	sl.SessionService = services.InitSessionService(&sl, "_session", 3600, logger)

	sl.ChatService = services.InitChatService(&sl, queries, logger)
	// initialize the rest of our services + http server
	s := &Server{
		router:   r,
		config:   config,
		logger:   logger,
		services: &sl,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%s", config.GoPort),
			Handler: r,
		},
	}
	s.routes()

	return s
}

// example of the route closures
func (s *Server) handleAnyHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))

}

func (s *Server) Open(port string) (err error) {
	s.logger.Info("Starter started!")
	return s.server.ListenAndServe()

}

func (s *Server) Close() error {

	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)

	defer cancel()
	return s.server.Shutdown(ctx)

}

// safe csrf getting
func csrfFromRequest(r *http.Request) string {
	return csrf.Token(r)

}
