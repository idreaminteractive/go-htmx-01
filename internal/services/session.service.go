package services

import (
	"fmt"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

	"github.com/sirupsen/logrus"
)

type ISessionService interface {
	WriteSession(c echo.Context, sp SessionPayload) error
	ReadSession(c echo.Context) (SessionPayload, error)
}

type SessionPayload struct {
	Email  string `json:"email"`
	UserId int    `json:"userId"`
}

type SessionService struct {
	sessionName string
	maxAge      int
	sl          *ServiceLocator
}

func InitSessionService(sl *ServiceLocator, sessionName string, maxAge int) *SessionService {
	return &SessionService{
		sessionName: sessionName,
		maxAge:      maxAge,
		sl:          sl,
	}
}

func (ss *SessionService) ReadSession(r *http.Request) (SessionPayload, error) {
	// this feels janky - but it's fine for now.
	sess, err := session.Get(ss.sessionName, c)
	if err != nil {
		logrus.Error("Error in getting session")
		return SessionPayload{}, err
	}
	payload := sess.Values["data"]
	if payload == nil {
		return SessionPayload{}, nil
	}
	return payload.(SessionPayload), nil

}



func (ss *SessionService) WriteSession(w http.ResponseWriter, r *http.Request, sp SessionPayload) error {
	sess, err := session.Get(ss.sessionName, c)
	if err != nil {
		logrus.Info("Could not get session")
		return err
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   ss.maxAge,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	// Set user as authenticated
	sess.Values["data"] = sp
	if err = sess.Save(r, w); err != nil {
		logrus.WithField("error", err).Error("Error in saving session")
		return err
	}

	return nil

}





type (
	// Config defines the config for Session middleware.
	Config struct {	
		// Session store.
		// Required.
		Store sessions.Store
	}
)

const (
	key = "_session_store"
)

var (
	// DefaultConfig is the default Session middleware config.
	DefaultConfig = Config{
		
	}
)

// Get returns a named session.
func Get(name string, r *http.Request) (*sessions.Session, error) {
	s := r.Context().Value(key)
	if s == nil {
		return nil, fmt.Errorf("%q session store not found", key)
	}
	store := s.(sessions.Store)
	return store.Get(r, name)
}

// Middleware returns a Session middleware.
func Middleware(store sessions.Store) func(next http.Handler) http.Handler {
	c := DefaultConfig
	c.Store = store
	return MiddlewareWithConfig(c)
}

// MiddlewareWithConfig returns a Sessions middleware with config.
// See `Middleware()`.
func MiddlewareWithConfig(config Config)  func(next http.Handler) http.Handler {
	// Defaults	
	if config.Store == nil {
		panic("echo: session middleware requires store")
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer context.Clear(r)
			ctx := context. .WithValue(r.Context(), key, config.Store)
    		next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
