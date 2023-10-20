package http

import (
	"main/internal/db"
	"main/internal/services"
	"main/internal/views/dto"
	"main/tests"
	"strings"

	"net/http"
	"net/http/httptest"
	"net/url"

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

func (m *MockAuthService) Authenticate(payload dto.UserLoginDTO) (*db.User, error) {
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
	e := setupEcho(EchoSetupStruct{DisableCSRF: true})
	e.GET("/login/", s.server.handleLoginGet)

	req := httptest.NewRequest(http.MethodGet, "/login/", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)
	assert.Equal(s.T(), http.StatusOK, rec.Code, "Bad status")

}

func (s *ServiceTestSuite) TestPostLogin_MissingFields() {
	e := setupEcho(EchoSetupStruct{DisableCSRF: true})
	e.POST("/login/", s.server.handleLoginPost)
	f := make(url.Values)
	f.Set("csrf", "stuff")
	f.Set("password", faker.Password())
	// f.Set("email", "not an email")
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)
	assert.Equal(s.T(), http.StatusBadRequest, rec.Code)

}

func (s *ServiceTestSuite) TestPostLogin_InvalidData() {
	e := setupEcho(EchoSetupStruct{DisableCSRF: true})
	e.POST("/login/", s.server.handleLoginPost)
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
	e := setupEcho(EchoSetupStruct{DisableCSRF: true})
	e.POST("/login/", s.server.handleLoginPost)
	f := make(url.Values)
	f.Set("password", faker.Password())
	f.Set("email", faker.Email())
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(f.Encode()))

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusOK, rec.Code)

}

func TestHttpPublicSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
