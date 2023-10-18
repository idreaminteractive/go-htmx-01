package http

import (
	"main/internal/db"
	"main/tests"
	"strings"

	"net/http"
	"net/http/httptest"
	"net/url"

	"main/internal/views"

	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"

	"github.com/stretchr/testify/suite"
)

// we need to mock our auth service
type ServiceTestSuite struct {
	suite.Suite
	server *Server
}

type MockAuthService struct {
}

func (m *MockAuthService) Authenticate(payload views.UserLoginDTO) (*db.User, error) {
	return &db.User{}, nil
}

// todo - refactor this out to make a global setup + tear down...
func (s *ServiceTestSuite) SetupSuite() {
	tests.GeneralSuiteSetup()

}

func (s *ServiceTestSuite) SetupTest() {
	// hook it up
	s.server = &Server{authenticationService: &MockAuthService{}}
}

func (s *ServiceTestSuite) TestGetLogin() {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/login")
	// the main thing here is testing the controllers
	// we can unit test the components in templ in another spot under views
	if s.NoError(s.server.handleLoginGet(c)) {
		s.Equal(http.StatusOK, rec.Code)

	}

}

func (s *ServiceTestSuite) TestPostLogin_MissingFields() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	f := make(url.Values)
	f.Set("password", faker.Password())
	// f.Set("email", "not an email")
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// the main thing here is testing the controllers
	// we can unit test the components in templ in another spot under views
	if s.NoError(s.server.handleLoginPost(c)) {
		s.Equal(http.StatusBadRequest, rec.Code)

	}

}

func (s *ServiceTestSuite) TestPostLogin_InvalidData() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	f := make(url.Values)
	f.Set("password", faker.Password())
	f.Set("email", "not an email")
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if s.NoError(s.server.handleLoginPost(c)) {
		s.Equal(http.StatusBadRequest, rec.Code)

	}

}

func (s *ServiceTestSuite) TestPostLogin_HappyPath() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	f := make(url.Values)
	f.Set("password", faker.Password())
	f.Set("email", faker.Email())
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if s.NoError(s.server.handleLoginPost(c)) {
		s.Equal(http.StatusOK, rec.Code)

	}

}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
