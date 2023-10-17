package services

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// todo - test!
// setup
// create random db filename
// setup db
// create queries obj
// create service + instantiate

type ServiceTestSuite struct {
	suite.Suite
}

func (s *ServiceTestSuite) SetupSuite() {
	// create the db
	s.T().Log("SetupSuite")
}

func (s *ServiceTestSuite) SetupTest() {
	// setup services + values
	s.T().Log("SetupTest")
}

func (s *ServiceTestSuite) BeforeTest(suiteName, testName string) {
	// clear + reboot the db
	// setup mocks, etc
	s.T().Logf("BeforeTest, %v, %v\n", suiteName, testName)
}

func (s *ServiceTestSuite) Test() {
	s.T().Log("Test")
}

func (s *ServiceTestSuite) AfterTest(suiteName, testName string) {
	// clear mocks
	s.T().Logf("AfterTest, %v, %v\n", suiteName, testName)
}

func (s *ServiceTestSuite) TearDownTest() {
	// nothing?
	s.T().Log("TearDownTest")
}

func (s *ServiceTestSuite) TearDownSuite() {
	s.T().Log("TearDownSuite")
	// remove the db
}

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
