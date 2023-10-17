package services

import (
	"context"

	"main/internal/db"
	"main/internal/views"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type IAuthenticationService interface {
	Authenticate(payload views.UserLoginDTO) (*db.User, error)
}

type AuthenticationService struct {
	Queries *db.Queries
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (as *AuthenticationService) Authenticate(payload views.UserLoginDTO) (*db.User, error) {
	ctx := context.Background()
	logrus.WithField("user", payload.Email).Info("Auth attempt")
	results, err := as.Queries.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		// since this is return one or fail - if it errors, the suer is not made.
		logrus.Errorf("No user found for %s, making new acct", payload.Email)
		// hash our password
		hashed, err := hashPassword(payload.Password)
		if err != nil {
			logrus.WithError(err).Error("Error in hashing password")
			return nil, &Error{Code: EINTERNAL, Message: "Could not hash password"}
		}
		// ok - hashing is cool.
		// make the user
		createdUser, err := as.Queries.CreateUser(ctx, db.CreateUserParams{Email: payload.Email, Password: hashed})
		if err != nil {
			logrus.WithError(err).Error("Error in creating user")
			// need to return proper errors
			return nil, &Error{Code: EINTERNAL, Message: "Could not create user"}
		}

		return &createdUser, nil
	}
	// user exists...
	if checkPasswordHash(payload.Password, results.Password) {
		return &results, nil
	} else {
		// this needs to return a new auth error
		return nil, &Error{Code: EUNAUTHORIZED, Message: "Invalid password"}
	}

}
