package services

import (
	"context"
	"fmt"
	"main/internal/db"
	"main/internal/sqlite"
	"main/internal/views/dto"
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

	sl := ServiceLocator{}
	queries := db.New(s.database.Connection)
	s.authService = &AuthenticationService{sl: &sl, queries: queries}
}

func (s *ServiceTestSuite) BeforeTest(suiteName, testName string) {

}

func (s *ServiceTestSuite) TestRegistrationSuccess() {
	users, err := s.authService.queries.GetAllUsers(s.context)
	s.Nil(err, "Error is not nil in getting users")
	s.Equal(len(users), 0)
	email, password := faker.Email(), faker.Password()
	u, err := s.authService.Register(dto.RegisterDTO{Handle: faker.Username(), Email: email, Password: password, ConfirmPassword: password})
	s.Nil(err, "Error is not nil")
	s.NotNil(u, "User does not exist")
	s.Equal(u.Email, email)
	// since it's hashed
	s.NotEqual(u.Password, password)
	users, err = s.authService.queries.GetAllUsers(s.context)
	s.Nil(err, "Error is not nil in getting users")
	s.Equal(1, len(users), "user was not created?")
}

func (s *ServiceTestSuite) TestRegistrationFailure() {
	users, err := s.authService.queries.GetAllUsers(s.context)
	s.Nil(err, "Error is not nil in getting users")
	s.Equal(len(users), 0)

	// create a user then try and create another with same email
	email, password := faker.Email(), faker.Password()
	u, err := s.authService.Register(dto.RegisterDTO{Handle: faker.Username(), Email: email, Password: password, ConfirmPassword: password})
	s.Nil(err, "Error should  be nil")
	s.NotNil(u, "User should not be nil")

	u, err = s.authService.Register(dto.RegisterDTO{Handle: faker.Username(), Email: email, Password: password, ConfirmPassword: "confirm"})
	s.NotNil(err, "Error is nil")
	s.Nil(u, "User not created")

	users, err = s.authService.queries.GetAllUsers(s.context)
	s.Nil(err, "Error is not nil in getting users")
	// Just the one
	s.Equal(1, len(users), "Only 1 user should be created")

}

func (s *ServiceTestSuite) TestAuthenticateSuccessfulLogin() {
	users, err := s.authService.queries.GetAllUsers(s.context)
	s.Nil(err, "Error is not nil in getting users")
	s.Equal(len(users), 0)
	email, password := faker.Email(), faker.Password()

	// first need to register them
	u, err := s.authService.Register(dto.RegisterDTO{Handle: faker.Username(), Email: email, Password: password, ConfirmPassword: password})
	s.Nil(err, "Error should  be nil")
	s.NotNil(u, "User should not be nil")

	u, err = s.authService.Authenticate(dto.UserLoginDTO{Email: email, Password: password})
	s.Nil(err, "Error is not nil")
	s.NotNil(u, "User does not exist")
	s.Equal(u.Email, email)
	// since it's hashed
	s.NotEqual(u.Password, password)
	users, err = s.authService.queries.GetAllUsers(s.context)
	s.Nil(err, "Error is not nil in getting users")
	s.Equal(1, len(users), "user was not created?")

}
func (s *ServiceTestSuite) TestBadPassword() {
	email, password, badPassword := faker.Email(), faker.Password(), faker.Password()
	// first need to register them
	u, err := s.authService.Register(dto.RegisterDTO{Handle: faker.Username(), Email: email, Password: password, ConfirmPassword: password})
	s.Nil(err, "Error should  be nil")
	s.NotNil(u, "User should not be nil")

	_, err = s.authService.Authenticate(dto.UserLoginDTO{Email: email, Password: badPassword})
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
