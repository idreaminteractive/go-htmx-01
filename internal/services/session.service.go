package services

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

	"github.com/gorilla/sessions"
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

func (ss *SessionService) ReadSession(c echo.Context) (SessionPayload, error) {
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

func (ss *SessionService) WriteSession(c echo.Context, sp SessionPayload) error {
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
	if err = sess.Save(c.Request(), c.Response()); err != nil {
		logrus.WithField("error", err).Error("Error in saving session")
		return err
	}

	return nil

}
