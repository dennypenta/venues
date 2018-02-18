package controllers

import (
	"venues/cmd/models"

	"net/http/httptest"
	"venues/cmd/fixtures"

	"errors"
	"net/http"
	"testing"

	"encoding/json"

	"strings"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gopkg.in/mgo.v2"
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

func (m *MockRepo) List(filter *models.Restaurant) ([]models.Restaurant, error) {
	args := m.Called(filter)
	return args.Get(0).([]models.Restaurant), args.Error(1)
}

func (m *MockRepo) Create(object *models.Restaurant) error {
	args := m.Called(object)
	return args.Error(0)
}

func (m *MockRepo) Update(query *models.Restaurant, object *models.Restaurant) error {
	args := m.Called(query, object)
	return args.Error(0)
}

func (m *MockRepo) Remove(query *models.Restaurant) error {
	args := m.Called(query)
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
		mockRepo.On(
			"List",
			mock.MatchedBy(func(i *models.Restaurant) bool {return true}),
		).Return(returnValue, nil)
		mockBinder := &MockBinder{}
		mockBinder.On(
			"Bind",
			mock.MatchedBy(func(i interface{}) bool { return true }),
			suite.echoContext,
		).Return(nil)

		suite.echoContext.Echo().Binder = mockBinder

		suite.controller.List(suite.echoContext)

		var resultValue []models.Restaurant
		json.NewDecoder(suite.recorder.Body).Decode(&resultValue)

		mockRepo.AssertExpectations(suite.T())
		mockBinder.AssertExpectations(suite.T())
		suite.Assertions.Equal(resultValue, returnValue)
		suite.Assertions.Equal(suite.echoContext.Response().Status, http.StatusOK)
	}
}

func (suite *RestaurantControllerTestSuite) TestListFailService() {
	var returnValue []models.Restaurant

	mockRepo := &MockRepo{}
	suite.controller = &RestaurantController{Repo: mockRepo}
	mockRepo.On(
		"List",
		mock.MatchedBy(func(i *models.Restaurant) bool {return true}),
	).Return(returnValue, errors.New("mocked error"))
	mockBinder := &MockBinder{}
	mockBinder.On(
		"Bind",
		mock.MatchedBy(func(i interface{}) bool { return true }),
		suite.echoContext,
	).Return(nil)

	suite.echoContext.Echo().Binder = mockBinder

	suite.controller.List(suite.echoContext)

	mockRepo.AssertExpectations(suite.T())
	mockBinder.AssertExpectations(suite.T())
	suite.Assertions.Equal(suite.echoContext.Response().Status, http.StatusServiceUnavailable)
}

func (suite *RestaurantControllerTestSuite) TestListFailBind() {
	mockBinder := &MockBinder{}
	mockBinder.On(
		"Bind",
		mock.MatchedBy(func(i interface{}) bool { return true }),
		suite.echoContext,
	).Return(errors.New("bind error"))

	suite.echoContext.Echo().Binder = mockBinder

	suite.controller.List(suite.echoContext)

	mockBinder.AssertExpectations(suite.T())
	suite.Assertions.Equal(suite.echoContext.Response().Status, http.StatusBadRequest)
}

func (suite *RestaurantControllerTestSuite) TestCreateSuccess() {
	body := "body"
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(body))
	suite.echoContext = echo.New().NewContext(req, suite.recorder)

	mockRepo := &MockRepo{}
	mockBinder := &MockBinder{}
	suite.controller = &RestaurantController{Repo: mockRepo}
	mockRepo.On(
		"Create",
		mock.MatchedBy(func(obj *models.Restaurant) bool { return true }),
	).Return(nil)
	mockBinder.On(
		"Bind",
		mock.MatchedBy(func(i interface{}) bool { return true }),
		suite.echoContext,
	).Return(nil)

	suite.echoContext.Echo().Binder = mockBinder

	suite.controller.Create(suite.echoContext)

	mockBinder.AssertExpectations(suite.T())
	mockRepo.AssertExpectations(suite.T())
	suite.Assertions.Equal(suite.echoContext.Response().Status, http.StatusOK)
}

func (suite *RestaurantControllerTestSuite) TestCreateFailFromRepo() {
	body := "body"
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(body))
	suite.echoContext = echo.New().NewContext(req, suite.recorder)

	mockRepo := &MockRepo{}
	mockBinder := &MockBinder{}
	suite.controller = &RestaurantController{Repo: mockRepo}
	mockRepo.On(
		"Create",
		mock.MatchedBy(func(obj *models.Restaurant) bool { return true }),
	).Return(errors.New("create error"))
	mockBinder.On(
		"Bind",
		mock.MatchedBy(func(i interface{}) bool { return true }),
		suite.echoContext,
	).Return(nil)

	suite.echoContext.Echo().Binder = mockBinder

	suite.controller.Create(suite.echoContext)

	mockBinder.AssertExpectations(suite.T())
	mockRepo.AssertExpectations(suite.T())
	suite.Assertions.Equal(suite.echoContext.Response().Status, http.StatusServiceUnavailable)
}

func (suite *RestaurantControllerTestSuite) TestCreateFailFromBind() {
	body := "body"
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(body))
	suite.echoContext = echo.New().NewContext(req, suite.recorder)

	mockBinder := &MockBinder{}
	suite.controller = &RestaurantController{}
	mockBinder.On(
		"Bind",
		mock.MatchedBy(func(i interface{}) bool { return true }),
		suite.echoContext,
	).Return(errors.New("bind error"))

	suite.echoContext.Echo().Binder = mockBinder

	suite.controller.Create(suite.echoContext)

	mockBinder.AssertExpectations(suite.T())

	suite.Assertions.Equal(suite.echoContext.Response().Status, http.StatusBadRequest)
}

func (suite *RestaurantControllerTestSuite) TestUpdateSuccess() {
	body := "body"
	req := httptest.NewRequest(echo.PATCH, "/", strings.NewReader(body))
	suite.echoContext = echo.New().NewContext(req, suite.recorder)
	suite.echoContext.SetParamNames("restaurant_id")
	suite.echoContext.SetParamValues("123")

	mockRepo := &MockRepo{}
	mockRepo.On(
		"Update",
		mock.MatchedBy(func(obj *models.Restaurant) bool { return true }),
		mock.MatchedBy(func(obj *models.Restaurant) bool { return true }),
	).Return(nil)
	suite.controller = &RestaurantController{Repo: mockRepo}
	mockBinder := &MockBinder{}
	mockBinder.On(
		"Bind",
		mock.MatchedBy(func(i interface{}) bool { return true }),
		suite.echoContext,
	).Return(nil)
	suite.echoContext.Echo().Binder = mockBinder

	suite.controller.Update(suite.echoContext)

	mockRepo.AssertExpectations(suite.T())
	mockBinder.AssertExpectations(suite.T())
	suite.Assertions.Equal(suite.echoContext.Response().Status, http.StatusOK)
}

func (suite *RestaurantControllerTestSuite) TestUpdateFailNotFound() {
	body := "body"
	req := httptest.NewRequest(echo.PATCH, "/", strings.NewReader(body))
	suite.echoContext = echo.New().NewContext(req, suite.recorder)
	suite.echoContext.SetParamNames("restaurant_id")
	suite.echoContext.SetParamValues("123")

	mockRepo := &MockRepo{}
	mockRepo.On(
		"Update",
		mock.MatchedBy(func(obj *models.Restaurant) bool { return true }),
		mock.MatchedBy(func(obj *models.Restaurant) bool { return true }),
	).Return(mgo.ErrNotFound)
	suite.controller = &RestaurantController{Repo: mockRepo}
	mockBinder := &MockBinder{}
	mockBinder.On(
		"Bind",
		mock.MatchedBy(func(i interface{}) bool { return true }),
		suite.echoContext,
	).Return(nil)
	suite.echoContext.Echo().Binder = mockBinder

	suite.controller.Update(suite.echoContext)

	mockRepo.AssertExpectations(suite.T())
	mockBinder.AssertExpectations(suite.T())
	suite.Assertions.Equal(suite.echoContext.Response().Status, http.StatusNotFound)
}

func (suite *RestaurantControllerTestSuite) TestUpdateFailFromBind() {
	body := "body"
	req := httptest.NewRequest(echo.PATCH, "/", strings.NewReader(body))
	suite.echoContext = echo.New().NewContext(req, suite.recorder)
	suite.echoContext.SetParamNames("restaurant_id")
	suite.echoContext.SetParamValues("123")

	suite.controller = &RestaurantController{}
	mockBinder := &MockBinder{}
	mockBinder.On(
		"Bind",
		mock.MatchedBy(func(i interface{}) bool { return true }),
		suite.echoContext,
	).Return(errors.New("bind error"))
	suite.echoContext.Echo().Binder = mockBinder

	suite.controller.Update(suite.echoContext)

	mockBinder.AssertExpectations(suite.T())
	suite.Assertions.Equal(suite.echoContext.Response().Status, http.StatusBadRequest)
}

func (suite *RestaurantControllerTestSuite) TestUpdateFailFromRepo() {
	body := "body"
	req := httptest.NewRequest(echo.PATCH, "/", strings.NewReader(body))
	suite.echoContext = echo.New().NewContext(req, suite.recorder)
	suite.echoContext.SetParamNames("restaurant_id")
	suite.echoContext.SetParamValues("123")

	mockRepo := &MockRepo{}
	mockRepo.On(
		"Update",
		mock.MatchedBy(func(obj *models.Restaurant) bool { return true }),
		mock.MatchedBy(func(obj *models.Restaurant) bool { return true }),
	).Return(errors.New("repo error"))
	suite.controller = &RestaurantController{Repo: mockRepo}
	mockBinder := &MockBinder{}
	mockBinder.On(
		"Bind",
		mock.MatchedBy(func(i interface{}) bool { return true }),
		suite.echoContext,
	).Return(nil)
	suite.echoContext.Echo().Binder = mockBinder

	suite.controller.Update(suite.echoContext)

	mockRepo.AssertExpectations(suite.T())
	mockBinder.AssertExpectations(suite.T())
	suite.Assertions.Equal(suite.echoContext.Response().Status, http.StatusServiceUnavailable)
}

func (suite *RestaurantControllerTestSuite) TestRemoveSuccess() {
	req := httptest.NewRequest(echo.DELETE, "/", nil)
	suite.echoContext = echo.New().NewContext(req, suite.recorder)
	suite.echoContext.SetParamNames("restaurant_id")
	suite.echoContext.SetParamValues("123")

	mockRepo := &MockRepo{}
	mockRepo.On(
		"Remove",
		mock.MatchedBy(func(obj *models.Restaurant) bool { return true }),
	).Return(nil)
	suite.controller = &RestaurantController{Repo: mockRepo}

	suite.controller.Remove(suite.echoContext)

	mockRepo.AssertExpectations(suite.T())
	suite.Assertions.Equal(suite.echoContext.Response().Status, http.StatusOK)
}

func (suite *RestaurantControllerTestSuite) TestRemoveFailNotFound() {
	req := httptest.NewRequest(echo.DELETE, "/", nil)
	suite.echoContext = echo.New().NewContext(req, suite.recorder)
	suite.echoContext.SetParamNames("restaurant_id")
	suite.echoContext.SetParamValues("123")

	mockRepo := &MockRepo{}
	mockRepo.On(
		"Remove",
		mock.MatchedBy(func(obj *models.Restaurant) bool { return true }),
	).Return(mgo.ErrNotFound)
	suite.controller = &RestaurantController{Repo: mockRepo}

	suite.controller.Remove(suite.echoContext)

	mockRepo.AssertExpectations(suite.T())
	suite.Assertions.Equal(suite.echoContext.Response().Status, http.StatusNotFound)
}

func (suite *RestaurantControllerTestSuite) TestRemoveFailService() {
	req := httptest.NewRequest(echo.DELETE, "/", nil)
	suite.echoContext = echo.New().NewContext(req, suite.recorder)
	suite.echoContext.SetParamNames("restaurant_id")
	suite.echoContext.SetParamValues("123")

	mockRepo := &MockRepo{}
	mockRepo.On(
		"Remove",
		mock.MatchedBy(func(obj *models.Restaurant) bool { return true }),
	).Return(errors.New("mocked error"))
	suite.controller = &RestaurantController{Repo: mockRepo}

	suite.controller.Remove(suite.echoContext)

	mockRepo.AssertExpectations(suite.T())
	suite.Assertions.Equal(suite.echoContext.Response().Status, http.StatusServiceUnavailable)
}

func TestRestaurantControllerTestSuite(t *testing.T) {
	suite.Run(t, new(RestaurantControllerTestSuite))
}
