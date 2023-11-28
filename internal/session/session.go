package session

import (
	"fmt"
	"main/internal/helpers"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
)

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
	DefaultConfig = Config{}
)

// Get returns a named session.
func Get(name string, r *http.Request) (*sessions.Session, error) {
	val := r.Context().Value(key)
	if val == nil {
		return nil, fmt.Errorf("%q session store not found", key)
	}
	store := val.(sessions.Store)
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
func MiddlewareWithConfig(config Config) func(next http.Handler) http.Handler {
	// Defaults

	if config.Store == nil {
		panic("server: session middleware requires store")
	}

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			defer context.Clear(r)
			cr := helpers.ContextSave(r, key, config.Store)
			next.ServeHTTP(w, cr)

		})
	}

	// return func(next echo.HandlerFunc) echo.HandlerFunc {
	// 	return func(c echo.Context) error {
	// 		if config.Skipper(c) {
	// 			return next(c)
	// 		}
	// 		defer context.Clear(c.Request())
	// 		c.Set(key, config.Store)
	// 		return next(c)
	// 	}
	// }
}
