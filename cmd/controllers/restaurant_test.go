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
	"strings"
)

type MockBinder struct {
	mock.Mock
}

func (m *MockBinder) Bind(i interface{}, c echo.Context) error {
	args := m.Called(i, c)
	return args.Error(0)
}


type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) List() ([]models.Restaurant, error) {
	args := m.Called()
	return args.Get(0).([]models.Restaurant), args.Error(1)
}

func (m *MockRepo) Create(object *models.Restaurant) error {
	args := m.Called(object)
	return args.Error(0)
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

func (suite *RestaurantControllerTestSuite) TestListSuccess() {
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

func (suite *RestaurantControllerTestSuite) TestListFail() {
	var returnValue []models.Restaurant

	mockRepo := &MockRepo{}
	suite.controller = &RestaurantController{Repo: mockRepo}
	mockRepo.On("List").Return(returnValue, errors.New("mocked error"))

	suite.controller.List(suite.echoContext)

	suite.Assertions.Equal(suite.echoContext.Response().Status, http.StatusServiceUnavailable)
}

func (suite *RestaurantControllerTestSuite) TestCreateSuccess() {
	body := "body"
	req := httptest.NewRequest(echo.GET, "/", strings.NewReader(body))
	suite.echoContext = echo.New().NewContext(req, suite.recorder)

	mockRepo := &MockRepo{}
	mockBinder := &MockBinder{}
	suite.controller = &RestaurantController{Repo: mockRepo}
	mockRepo.On(
	"Create",
		mock.MatchedBy(func(obj *models.Restaurant) bool {return true}),
	).Return(nil)
	mockBinder.On(
	"Bind",
		mock.MatchedBy(func(i interface{}) bool {return true}),
		suite.echoContext,
	).Return(nil)

	suite.echoContext.Echo().Binder = mockBinder

	suite.controller.Create(suite.echoContext)

	mockBinder.AssertExpectations(suite.T())
	mockRepo.AssertExpectations(suite.T())
	suite.Assertions.Equal(suite.echoContext.Response().Status, http.StatusOK)
}

func (suite *RestaurantControllerTestSuite) TestCreateFailByRepo() {
	body := "body"
	req := httptest.NewRequest(echo.GET, "/", strings.NewReader(body))
	suite.echoContext = echo.New().NewContext(req, suite.recorder)

	mockRepo := &MockRepo{}
	mockBinder := &MockBinder{}
	suite.controller = &RestaurantController{Repo: mockRepo}
	mockRepo.On(
		"Create",
		mock.MatchedBy(func(obj *models.Restaurant) bool {return true}),
	).Return(errors.New("create error"))
	mockBinder.On(
		"Bind",
		mock.MatchedBy(func(i interface{}) bool {return true}),
		suite.echoContext,
	).Return(nil)

	suite.echoContext.Echo().Binder = mockBinder

	suite.controller.Create(suite.echoContext)

	mockBinder.AssertExpectations(suite.T())
	mockRepo.AssertExpectations(suite.T())
	suite.Assertions.Equal(suite.echoContext.Response().Status, http.StatusServiceUnavailable)
}

func (suite *RestaurantControllerTestSuite) TestCreateFailByBind() {
	body := "body"
	req := httptest.NewRequest(echo.GET, "/", strings.NewReader(body))
	suite.echoContext = echo.New().NewContext(req, suite.recorder)

	mockBinder := &MockBinder{}
	suite.controller = &RestaurantController{}
	mockBinder.On(
		"Bind",
		mock.MatchedBy(func(i interface{}) bool {return true}),
		suite.echoContext,
	).Return(errors.New("bind error"))

	suite.echoContext.Echo().Binder = mockBinder

	suite.controller.Create(suite.echoContext)

	mockBinder.AssertExpectations(suite.T())

	suite.Assertions.Equal(suite.echoContext.Response().Status, http.StatusBadRequest)
}

func TestRestaurantControllerTestSuite(t *testing.T) {
	suite.Run(t, new(RestaurantControllerTestSuite))
}
