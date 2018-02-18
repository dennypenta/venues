package assembly

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"venues/pkg/healthcheckers"

	"github.com/labstack/echo"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) Check() error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockService) Message() string {
	return ""
}

type HealthCheckTestSuite struct {
	suite.Suite

	echoContext echo.Context
	healthCheck *HealthCheck
}

func (suite *HealthCheckTestSuite) SetupTest() {
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	suite.echoContext = echo.New().NewContext(req, rec)
}

func (suite *HealthCheckTestSuite) TestSuccess() {
	mockService := &MockService{}
	suite.healthCheck = &HealthCheck{[]healthcheckers.Checker{mockService}}

	mockService.On("Check").Return(nil)

	suite.healthCheck.Check(suite.echoContext)

	mockService.AssertExpectations(suite.T())
	suite.Assertions.Equal(suite.echoContext.Response().Status, http.StatusOK)
}

func (suite *HealthCheckTestSuite) TestFail() {
	mockService := &MockService{}
	suite.healthCheck = &HealthCheck{[]healthcheckers.Checker{mockService}}

	mockService.On("Check").Return(errors.New("mocked error"))

	suite.healthCheck.Check(suite.echoContext)

	mockService.AssertExpectations(suite.T())
	suite.Assertions.Equal(suite.echoContext.Response().Status, http.StatusServiceUnavailable)
}

func TestHealthCheck(t *testing.T) {
	suite.Run(t, new(HealthCheckTestSuite))
}
