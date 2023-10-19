package services

import (
	"context"
	"fmt"
	"main/internal/db"
	"main/internal/sqlite"
	"main/internal/views"
	"main/tests"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"github.com/go-faker/faker/v4"

	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
	authService  *AuthenticationService
	context      context.Context
	database     *sqlite.DB
	databasePath string
}

// todo - refactor this out to make a global setup + tear down...
func (s *ServiceTestSuite) SetupSuite() {
	tests.GeneralSuiteSetup()
	s.context = context.Background()
	s.databasePath = fmt.Sprintf("./%s.db", faker.UUIDDigit())
	// create the db + run goose migrations
	err := tests.SetupTestDatabase(s.databasePath)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *ServiceTestSuite) SetupTest() {
	// hook it up
	s.database = sqlite.NewDB(s.databasePath)
	if err := s.database.Open(); err != nil {
		s.T().Fatal(err)
	}

	queries := db.New(s.database.Connection)
	s.authService = &AuthenticationService{Queries: queries}
}

func (s *ServiceTestSuite) BeforeTest(suiteName, testName string) {

}

func (s *ServiceTestSuite) TestAuthenticateNewUser() {
	users, err := s.authService.Queries.GetAllUsers(s.context)
	s.Nil(err, "Error is not nil in getting users")
	s.Equal(len(users), 0)
	email, password := faker.Email(), faker.Password()
	u, err := s.authService.Authenticate(views.UserLoginDTO{Email: email, Password: password})
	s.Nil(err, "Error is not nil")
	s.NotNil(u, "User does not exist")
	s.Equal(u.Email, email)
	// since it's hashed
	s.NotEqual(u.Password, password)
	users, err = s.authService.Queries.GetAllUsers(s.context)
	s.Nil(err, "Error is not nil in getting users")
	s.Equal(1, len(users), "user was not created?")

}

func (s *ServiceTestSuite) TestAuthenticateSuccessfulLogin() {
	// make our user first
	email, password := faker.Email(), faker.Password()
	u, err := s.authService.Authenticate(views.UserLoginDTO{Email: email, Password: password})
	s.Nil(err, "Error is not nil")
	s.NotNil(u, "User does not exist")
	s.Equal(u.Email, email)
	users, err := s.authService.Queries.GetAllUsers(s.context)
	s.Nil(err, "Error is not nil in getting users")
	s.Equal(len(users), 1)

	u, err = s.authService.Authenticate(views.UserLoginDTO{Email: email, Password: password})
	s.Nil(err, "Error is not nil")
	s.NotNil(u, "User does not exist")
	s.Equal(u.Email, email)

	users, err = s.authService.Queries.GetAllUsers(s.context)
	s.Nil(err, "Error is not nil in getting users")
	s.Equal(len(users), 1, "new user was created unintentionally")
}

func (s *ServiceTestSuite) TestBadPassword() {
	// make our user first
	email, password, badPassword := faker.Email(), faker.Password(), faker.Password()
	u, err := s.authService.Authenticate(views.UserLoginDTO{Email: email, Password: password})
	s.Nil(err, "Error is not nil")
	s.NotNil(u, "User does not exist")
	s.Equal(u.Email, email)

	_, err = s.authService.Authenticate(views.UserLoginDTO{Email: email, Password: badPassword})
	s.NotNil(err, "Error is nil")
}

// note - we don't test our controller in the service tests!

func (s *ServiceTestSuite) AfterTest(suiteName, testName string) {
	tests.WipeDB(s.database.Connection)
}

func (s *ServiceTestSuite) TearDownTest() {
	// nothing?

}

func (s *ServiceTestSuite) TearDownSuite() {
	// close db + cleanup
	s.database.Close()
	tests.TearDownSuite(s.databasePath)
}
func TestServicesAuthenticationSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
