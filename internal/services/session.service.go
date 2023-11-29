package services

import (
	"main/internal/session"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
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
}

func InitSessionService(sl *ServiceLocator, sessionName string, maxAge int) *SessionService {
	// var err error

	// secret, err := hex.DecodeString("13b4dff8f84a10851021ec8a5d12b570d562c92fe6b5ec4c4129f595bcb3234b")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	return &SessionService{
		sessionName: sessionName,
		maxAge:      maxAge,
		sl:          sl,
		// secret:      secret,
	}
}

func (ss *SessionService) ReadSession(r *http.Request) (SessionPayload, error) {
	// this feels janky - but it's fine for now.

	// sess, err := ss.getCookie(r)
	// if err != nil {
	// 	logrus.Error("Error in getting session")
	// 	return &SessionPayload{}, err
	// }

	sess, err := session.Get(ss.sessionName, r)
	if err != nil {
		logrus.Error("Error in getting session")
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
		logrus.Error(err)
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
	// if err := ss.setCookie(w, sp); err != nil {
	// 	logrus.WithError(err).Error("Could not write session")
	// 	return err
	// }

	// return nil

}

// func (ss *SessionService) setCookie(w http.ResponseWriter, sp SessionPayload) error {

// 	var buf bytes.Buffer

// 	err := gob.NewEncoder(&buf).Encode(sp)
// 	if err != nil {
// 		logrus.Error(err)
// 		return err
// 	}

// 	cookie := http.Cookie{
// 		Name:     ss.sessionName,
// 		Value:    buf.String(),
// 		Path:     "/",
// 		MaxAge:   ss.maxAge,
// 		HttpOnly: true,
// 		Secure:   true,
// 		SameSite: http.SameSiteLaxMode,
// 	}

// 	err = cookies.WriteEncrypted(w, cookie, ss.secret)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (ss *SessionService) getCookie(r *http.Request) (*SessionPayload, error) {
// 	gobEncodedValue, err := cookies.ReadEncrypted(r, ss.sessionName, ss.secret)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, http.ErrNoCookie):
// 			logrus.Error("Cookie not found")

// 		case errors.Is(err, cookies.ErrInvalidValue):
// 			logrus.Error("Invalid cookie")

// 		default:

// 			logrus.Error("server error")

// 		}
// 		return nil, err
// 	}

// 	var sp SessionPayload

// 	reader := strings.NewReader(gobEncodedValue)

// 	if err := gob.NewDecoder(reader).Decode(&sp); err != nil {
// 		logrus.Error(err)
// 		return nil, err
// 	}
// 	return &sp, nil
// }
