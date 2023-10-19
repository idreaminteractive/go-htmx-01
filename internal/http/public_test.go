package http

import (
	"main/internal/db"
	"main/internal/services"
	"main/tests"
	"strings"

	"net/http"
	"net/http/httptest"
	"net/url"

	"main/internal/views"

	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"

	"github.com/stretchr/testify/assert"
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

type MockSessionService struct {
}

func (mss *MockSessionService) WriteSession(c echo.Context, sp services.SessionPayload) error {
	return nil
}

func (mss *MockSessionService) ReadSession(c echo.Context) (services.SessionPayload, error) {
	return services.SessionPayload{}, nil
}

// todo - refactor this out to make a global setup + tear down...
func (s *ServiceTestSuite) SetupSuite() {
	tests.GeneralSuiteSetup()

}

func (s *ServiceTestSuite) SetupTest() {
	// hook it up
	s.server = &Server{
		authenticationService: &MockAuthService{},
		sessionService:        &MockSessionService{},
	}
}

// left off here. echo is setup + we need to setup routes!
func (s *ServiceTestSuite) TestGetLogin() {
	e := setupEcho(faker.Word())
	s.server.registerPublicRoutes()
	req := httptest.NewRequest(http.MethodGet, "/login/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(s.T(), http.StatusOK, rec.Code, "Bad status")

}

func (s *ServiceTestSuite) TestPostLogin_MissingFields() {
	e := setupEcho("")
	f := make(url.Values)
	f.Set("password", faker.Password())
	// f.Set("email", "not an email")
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)
	assert.Equal(s.T(), http.StatusBadRequest, rec.Code)

}

func (s *ServiceTestSuite) TestPostLogin_InvalidData() {
	e := setupEcho("")
	f := make(url.Values)
	f.Set("password", faker.Password())
	f.Set("email", "not an email")
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
}

func (s *ServiceTestSuite) TestPostLogin_HappyPath() {
	e := setupEcho("")
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

func TestHttpPublicSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
