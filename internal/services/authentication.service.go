package services

import (
	"context"
	"fmt"

	"main/internal/db"

	"main/internal/views/dto"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticationService struct {
	sl      *ServiceLocator
	queries *db.Queries
}

func InitAuthService(sl *ServiceLocator, queries *db.Queries) *AuthenticationService {
	return &AuthenticationService{
		sl:      sl,
		queries: queries,
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func checkPasswordHash(clearText, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(clearText))
	return err == nil
}

func (as *AuthenticationService) Register(payload dto.RegisterDTO) (*db.User, error) {
	ctx := context.Background()
	logrus.WithField("user", payload.Email).Info("Registration attempt")
	// first, hash our pw
	hashed, err := hashPassword(payload.Password)
	if err != nil {
		logrus.WithError(err).Error("Error in hashing password")
		return nil, &Error{Code: EINTERNAL, Message: "Could not hash password"}
	}

	// IF erro

	_, err = as.queries.GetUserByEmail(ctx, payload.Email)
	if err == nil {
		logrus.Error("User already exists!")
		// user exists already
		return nil, &Error{Code: ECONFLICT, Message: "User already exists"}
	}
	// user does not exist, create it and return it
	logrus.Errorf("No user found for %s, making new acct", payload.Email)

	// make the user
	fmt.Printf("Regging with %q", hashed)
	createdUser, err := as.queries.CreateUser(ctx, db.CreateUserParams{Email: payload.Email, Handle: payload.Handle, Password: hashed})
	if err != nil {
		logrus.WithError(err).Error("Error in creating user")
		// need to return proper errors
		return nil, &Error{Code: EINTERNAL, Message: "Could not create user"}
	}

	return &createdUser, nil

}

func (as *AuthenticationService) Authenticate(payload dto.UserLoginDTO) (*db.User, error) {
	ctx := context.Background()
	logrus.WithField("user", payload.Email).Info("Auth attempt")
	// first, hash our pw
	hashed, err := hashPassword(payload.Password)
	if err != nil {
		logrus.WithError(err).Error("Error in hashing password")
		return nil, &Error{Code: EINTERNAL, Message: "Could not hash password"}
	}

	results, err := as.queries.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		logrus.Errorf("No user found for %s, making new acct", payload.Email)

		return nil, &Error{Code: EUNAUTHORIZED, Message: "Invalid email or password"}
	}
	// user exists...
	if checkPasswordHash(payload.Password, results.Password) {
		return &results, nil
	} else {
		// this needs to return a new auth error
		fmt.Printf("%q -  %q - %q", payload.Password, hashed, results.Password)
		logrus.Error("Failed pass check")
		return nil, &Error{Code: EUNAUTHORIZED, Message: "Invalid email or password"}
	}

}
