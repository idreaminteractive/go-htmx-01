package services

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

type ISessionService interface {
	WriteSession(c echo.Context, sp *SessionPayload) error
	ReadSession(c echo.Context) (*SessionPayload, error)
}

type SessionPayload struct {
	UserId string `json:"userId"`
}

type SessionService struct {
	SessionName string
	MaxAge      int
}

func (ss *SessionService) ReadSession(c echo.Context) (*SessionPayload, error) {
	cookie, err := c.Cookie(cs.CookieName)
	if err != nil {
		return nil, err
	}

	// marshal the cookie value
	var sp SessionPayload
	if err := json.Unmarshal([]byte(cookie.Value), &cp); err != nil {
		return nil, err
	}
	// ok - we're good. return it!
	return &cp, nil
}

func (ss *SessionService) WriteSession(c echo.Context, sp SessionPayload) error {
	sess, err := session.Get(ss.SessionName, c)
	if err != nil {
		logrus.Info("Could not get session")
		return err
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   ss.MaxAge,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	// Set user as authenticated
	sess.Values["authenticated"] = true
	sess.Save(c.Request(), c.Response())
	cookie := new(http.Cookie)
	cookie.Name = cs.CookieName
	val, err := json.Marshal(SessionService{UserId: "dave", Csrf: "potato"})
	if err != nil {
		return err
	}
	// ok - so val is good + should be json byes.

	cookie.MaxAge = cs.MaxAge
	cookie.Value = string(val[:])
	c.SetCookie(cookie)
	return nil

}
