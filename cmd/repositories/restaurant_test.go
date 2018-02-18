package repositories

import (
	"errors"
	"testing"

	"venues/cmd/fixtures"
	"venues/cmd/models"
	"venues/cmd/settings"

	"venues/pkg/mongo"

	"venues/cmd/storages"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MockQuerier struct {
	mock.Mock
}

func (m *MockQuerier) All(model interface{}) error {
	args := m.Called(model)
	return args.Error(0)
}

func (m *MockQuerier) One(model interface{}) error {
	args := m.Called(model)
	return args.Error(0)
}

func (m *MockQuerier) matchedBy(value interface{}) bool {
	return true
}

type MockDataAccess struct {
	mock.Mock
}

func (m *MockDataAccess) Find(query interface{}) mongo.Querier {
	args := m.Called(query)
	return args.Get(0).(mongo.Querier)
}

func (m *MockDataAccess) Insert(model interface{}) error {
	args := m.Called(model)
	return args.Error(0)
}

func (m *MockDataAccess) Update(query interface{}, model interface{}) error {
	args := m.Called(query, model)
	return args.Error(0)
}

type RestaurantRepoTestSuite struct {
	suite.Suite

	storage *mgo.Collection
	repo    *RestaurantRepo
}

func (suite *RestaurantRepoTestSuite) SetupTest() {
	settings.Load()

	suite.storage = storages.GetTestStorage().C(models.RestaurantCollectionName)
	suite.repo = &RestaurantRepo{
		storage: &mongo.DataAccess{Collection: suite.storage},
	}
}

func (suite *RestaurantRepoTestSuite) TearDownTest() {
	suite.storage.DropCollection()
}

func (suite *RestaurantRepoTestSuite) TestSuccessList() {
	expected := fixtures.SimpleRestaurantSet()
	for _, i := range expected {
		if err := suite.storage.Insert(i); err != nil {
			suite.T().Fatal(err.Error())
		}
	}

	result, err := suite.repo.List()

	suite.Assertions.Nil(err)
	suite.Assertions.Equal(result, expected)
}

func (suite *RestaurantRepoTestSuite) TestEmptyList() {
	var expected []models.Restaurant

	result, err := suite.repo.List()

	suite.Assertions.Nil(err)
	suite.Assertions.Equal(result, expected)
}

func (suite *RestaurantRepoTestSuite) TestErrorList() {
	mockedStorage := &MockDataAccess{}
	mockedQuerier := &MockQuerier{}
	suite.repo = &RestaurantRepo{storage: mockedStorage}

	mockedStorage.On("Find", bson.M{}).Return(mockedQuerier)
	mockedQuerier.On("All", mock.MatchedBy(mockedQuerier.matchedBy)).Return(errors.New("Mocked error"))

	_, err := suite.repo.List()

	mockedStorage.AssertExpectations(suite.T())
	mockedQuerier.AssertExpectations(suite.T())
	suite.Assertions.Error(err)
}

func (suite *RestaurantRepoTestSuite) TestCreateSuccess() {
	data, _ := suite.repo.List()
	suite.Assertions.Zero(data)

	expected := &models.Restaurant{Name: "Name", Menu: []models.Dish{}}
	err := suite.repo.Create(expected)
	suite.Assertions.Nil(err)

	result := &models.Restaurant{}
	suite.repo.storage.Find(expected).One(result)
	suite.Assertions.Equal(result, expected)
}

func (suite *RestaurantRepoTestSuite) TestCreateError() {
	mockAccess := &MockDataAccess{}
	suite.repo = &RestaurantRepo{storage: mockAccess}

	object := &models.Restaurant{Name: "Name", Menu: []models.Dish{}}
	mockAccess.On("Insert", object).Return(errors.New("mocked error"))

	err := suite.repo.Create(object)

	mockAccess.AssertExpectations(suite.T())
	suite.Assertions.Error(err)
}

func (suite *RestaurantRepoTestSuite) TestUpdateSuccess() {
	object := &models.Restaurant{Name: "Name", Menu: []models.Dish{}}
	suite.repo.Create(object)

	update := &models.Restaurant{Name: "Name333", Menu: []models.Dish{}}
	err := suite.repo.Update(object, update)
	suite.Assertions.Nil(err)

	result := &models.Restaurant{}
	update.ID = object.ID
	suite.repo.storage.Find(update).One(result)
	suite.Assertions.Equal(result, update)
}

func (suite *RestaurantRepoTestSuite) TestUpdateError() {
	mockAccess := &MockDataAccess{}
	suite.repo = &RestaurantRepo{storage: mockAccess}

	object := &models.Restaurant{Name: "Name", Menu: []models.Dish{}}
	update := &models.Restaurant{Name: "Name333", Menu: []models.Dish{}}
	mockAccess.On("Update", object, update).Return(errors.New("mocked error"))

	err := suite.repo.Update(object, update)

	mockAccess.AssertExpectations(suite.T())
	suite.Assertions.Error(err)
}

func TestRestaurantRepoTestSuite(t *testing.T) {
	suite.Run(t, new(RestaurantRepoTestSuite))
}
