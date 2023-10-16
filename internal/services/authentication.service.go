package services

import (
	"context"
	"main/internal/db"
	"main/internal/views"

	"github.com/sirupsen/logrus"
)

type IAuthenticationService interface {
	Authenticate(payload views.UserLoginDTO) (*db.User, error)
}

type AuthenticationService struct {
	Queries *db.Queries
}

func (as *AuthenticationService) Authenticate(payload views.UserLoginDTO) (*db.User, error) {
	ctx := context.Background()
	logrus.WithField("user", payload.Email).Info("Auth attempt")
	results, err := as.Queries.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		// since this is return one or fail - if it errors, the suer is not made.
		logrus.Errorf("Login failure for %s", payload.Email)
		return nil, err
	}
	logrus.Infof("Login success for %s", payload.Email)
	return &results, nil

}
