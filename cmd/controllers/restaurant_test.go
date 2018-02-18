package controllers

import (
	"venues/cmd/models"

	"net/http/httptest"
	"venues/cmd/fixtures"

	"errors"
	"net/http"
	"testing"

	"encoding/json"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) List() ([]models.Restaurant, error) {
	args := m.Called()
	return args.Get(0).([]models.Restaurant), args.Error(1)
}

type RestaurantControllerTestSuite struct {
	suite.Suite

	controller  *RestaurantController
	echoContext echo.Context
	recorder    *httptest.ResponseRecorder
}

func (suite *RestaurantControllerTestSuite) SetupTest() {
	req := httptest.NewRequest(echo.GET, "/", nil)
	suite.recorder = httptest.NewRecorder()
	suite.echoContext = echo.New().NewContext(req, suite.recorder)

}

func (suite *RestaurantControllerTestSuite) TestSuccess() {
	var emptyData []models.Restaurant
	for _, returnValue := range [][]models.Restaurant{
		fixtures.SimpleRestaurantSet(),
		emptyData,
	} {
		mockRepo := &MockRepo{}
		suite.controller = &RestaurantController{Repo: mockRepo}
		mockRepo.On("List").Return(returnValue, nil)

		suite.controller.List(suite.echoContext)

		var resultValue []models.Restaurant
		json.NewDecoder(suite.recorder.Body).Decode(&resultValue)

		suite.Assertions.Equal(resultValue, returnValue)
		suite.Assertions.Equal(suite.echoContext.Response().Status, http.StatusOK)
	}
}

func (suite *RestaurantControllerTestSuite) TestFail() {
	var returnValue []models.Restaurant

	mockRepo := &MockRepo{}
	suite.controller = &RestaurantController{Repo: mockRepo}
	mockRepo.On("List").Return(returnValue, errors.New("mocked error"))

	suite.controller.List(suite.echoContext)

	suite.Assertions.Equal(suite.echoContext.Response().Status, http.StatusServiceUnavailable)
}

func TestRestaurantControllerTestSuite(t *testing.T) {
	suite.Run(t, new(RestaurantControllerTestSuite))
}
