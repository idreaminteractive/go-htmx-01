package services

import (
	"main/internal/session"
	"net/http"

	"github.com/go-chi/httplog/v2"
	"github.com/gorilla/sessions"
)

type ISessionService interface {
	WriteSession(w http.ResponseWriter, r *http.Request, sp SessionPayload) error
	ReadSession(r *http.Request) (SessionPayload, error)
}

type SessionPayload struct {
	Email  string `json:"email"`
	UserId int    `json:"userId"`
}

type SessionService struct {
	sessionName string
	maxAge      int
	sl          *ServiceLocator
	secret      []byte
	logger      *httplog.Logger
}

func InitSessionService(sl *ServiceLocator, sessionName string, maxAge int, logger *httplog.Logger) *SessionService {
	// var err error

	return &SessionService{
		sessionName: sessionName,
		maxAge:      maxAge,
		sl:          sl,
		logger:      logger,
	}
}

func (ss *SessionService) ReadSession(r *http.Request) (SessionPayload, error) {
	// this feels janky - but it's fine for now.

	sess, err := session.Get(ss.sessionName, r)
	if err != nil {
		ss.logger.Error("Error in getting session", err)
		return SessionPayload{}, err
	}
	payload := sess.Values["data"]
	if payload == nil {
		return SessionPayload{}, nil
	}
	return payload.(SessionPayload), nil

	// return sess, nil

}

func (ss *SessionService) WriteSession(w http.ResponseWriter, r *http.Request, sp SessionPayload) error {
	sess, err := session.Get(ss.sessionName, r)
	if err != nil {
		ss.logger.Error("Error getting session", err)
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
		ss.logger.Error("Error in saving session", err)
		return err
	}

	return nil

}
