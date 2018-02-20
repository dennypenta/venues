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

func (m *MockQuerier) Select(fields interface{}) mongo.Querier {
	args := m.Called(fields)
	return args.Get(0).(mongo.Querier)
}

func (m *MockQuerier) All(model interface{}) error {
	args := m.Called(model)
	return args.Error(0)
}

func (m *MockQuerier) One(model interface{}) error {
	args := m.Called(model)
	return args.Error(0)
}

func (m *MockQuerier) Sort(ordering string) mongo.Querier {
	args := m.Called(ordering)
	return args.Get(0).(mongo.Querier)
}

func (m *MockQuerier) Skip(n int) mongo.Querier {
	args := m.Called(n)
	return args.Get(0).(mongo.Querier)
}

func (m *MockQuerier) Limit(n int) mongo.Querier {
	args := m.Called(n)
	return args.Get(0).(mongo.Querier)
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

func (m *MockDataAccess) Remove(query interface{}) error {
	args := m.Called(query)
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

	result, err := suite.repo.List(&models.Restaurant{}, "", 0)

	suite.Assertions.Nil(err)
	suite.Assertions.Equal(result, expected)
}

func (suite *RestaurantRepoTestSuite) TestFilterList() {
	expected := fixtures.SimpleRestaurantSet()
	for _, i := range expected {
		if err := suite.storage.Insert(i); err != nil {
			suite.T().Fatal(err.Error())
		}
	}

	filter := &models.Restaurant{City: "City1"}
	result, err := suite.repo.List(filter, "", 0)

	suite.Assertions.Nil(err)
	for _, i := range result {
		suite.Assertions.Equal(i.City, "City1")
	}
}

func (suite *RestaurantRepoTestSuite) TestOrderingList() {
	expected := fixtures.SimpleRestaurantSet()
	for _, i := range expected {
		if err := suite.storage.Insert(i); err != nil {
			suite.T().Fatal(err.Error())
		}
	}

	ordering := "-rating"
	result, err := suite.repo.List(&models.Restaurant{}, ordering, 0)

	suite.Assertions.Nil(err)
	suite.Assertions.True(result[0].Rating > result[1].Rating)
}

func (suite *RestaurantRepoTestSuite) TestPaginateList() {
	expected := fixtures.SimpleRestaurantSet()
	for _, i := range expected {
		if err := suite.storage.Insert(i); err != nil {
			suite.T().Fatal(err.Error())
		}
	}

	result, err := suite.repo.List(&models.Restaurant{}, "", 2)

	suite.Assertions.Nil(err)
	suite.Assertions.True(len(result) == 0)
}

func (suite *RestaurantRepoTestSuite) TestEmptyList() {
	var expected []models.Restaurant

	result, err := suite.repo.List(&models.Restaurant{}, "", 0)

	suite.Assertions.Nil(err)
	suite.Assertions.Equal(result, expected)
}

func (suite *RestaurantRepoTestSuite) TestErrorList() {
	mockedStorage := &MockDataAccess{}
	mockedQuerier := &MockQuerier{}
	suite.repo = &RestaurantRepo{storage: mockedStorage}

	mockedStorage.On(
		"Find",
		mock.MatchedBy(func(i interface{}) bool { return true }),
	).Return(mockedQuerier)
	mockedQuerier.On(
		"Select",
		bson.M{"menu": 0},
	).Return(mockedQuerier)
	mockedQuerier.On(
		"All",
		mock.MatchedBy(func(i interface{}) bool { return true }),
	).Return(errors.New("mocked error"))

	_, err := suite.repo.List(&models.Restaurant{}, "", 0)

	mockedStorage.AssertExpectations(suite.T())
	mockedQuerier.AssertExpectations(suite.T())
	suite.Assertions.Error(err)
}

func (suite *RestaurantRepoTestSuite) TestCreateSuccess() {
	data, _ := suite.repo.List(&models.Restaurant{}, "", 0)
	suite.Assertions.Zero(data)

	expected := &models.Restaurant{Name: "Name"}
	err := suite.repo.Create(expected)
	suite.Assertions.Nil(err)

	result := &models.Restaurant{}
	suite.repo.storage.Find(expected).One(result)
	expected.ID = result.ID
	suite.Assertions.Equal(result, expected)
}

func (suite *RestaurantRepoTestSuite) TestCreateError() {
	mockAccess := &MockDataAccess{}
	suite.repo = &RestaurantRepo{storage: mockAccess}

	object := &models.Restaurant{Name: "Name"}
	mockAccess.On("Insert", object).Return(errors.New("mocked error"))

	err := suite.repo.Create(object)

	mockAccess.AssertExpectations(suite.T())
	suite.Assertions.Error(err)
}

func (suite *RestaurantRepoTestSuite) TestUpdateSuccess() {
	object := &models.Restaurant{Name: "Name"}
	suite.repo.Create(object)

	update := &models.Restaurant{City: "City33"}
	err := suite.repo.Update(object, update)
	suite.Assertions.Nil(err)

	result := &models.Restaurant{}
	suite.repo.storage.Find(update).One(result)
	suite.Assertions.Equal(result, &models.Restaurant{ID: result.ID, Name: object.Name, City: update.City})
}

func (suite *RestaurantRepoTestSuite) TestUpdateError() {
	mockAccess := &MockDataAccess{}
	suite.repo = &RestaurantRepo{storage: mockAccess}

	object := &models.Restaurant{Name: "Name"}
	update := &models.Restaurant{Name: "Name333"}
	mockAccess.On("Update", object, bson.M{"$set": update}).Return(errors.New("mocked error"))

	err := suite.repo.Update(object, update)

	mockAccess.AssertExpectations(suite.T())
	suite.Assertions.Error(err)
}

func (suite *RestaurantRepoTestSuite) TestRemoveSuccess() {
	data, _ := suite.repo.List(&models.Restaurant{}, "", 0)
	suite.Assertions.Zero(data)

	object := &models.Restaurant{Name: "Name"}
	err := suite.repo.Create(object)
	suite.Assertions.Nil(err)

	err = suite.repo.Remove(object)
	suite.Assertions.Nil(err)

	data, _ = suite.repo.List(&models.Restaurant{}, "", 0)
	suite.Assertions.Zero(data)
}

func (suite *RestaurantRepoTestSuite) TestRemoveError() {
	mockAccess := &MockDataAccess{}
	suite.repo = &RestaurantRepo{storage: mockAccess}

	object := &models.Restaurant{Name: "Name"}
	mockAccess.On("Remove", object).Return(errors.New("mocked error"))

	err := suite.repo.Remove(object)

	mockAccess.AssertExpectations(suite.T())
	suite.Assertions.Error(err)
}

func (suite *RestaurantRepoTestSuite) TestPushDishSuccess() {
	object := &models.Restaurant{Name: "Name"}
	suite.repo.Create(object)

	dish := &models.Dish{Name: "Name"}
	err := suite.repo.AddDish(object, dish)
	suite.Assertions.Nil(err)

	result := &models.Restaurant{}
	suite.repo.storage.Find(object).One(result)

	suite.Assertions.Equal(&result.Menu[0], dish)
}

func TestRestaurantRepoTestSuite(t *testing.T) {
	suite.Run(t, new(RestaurantRepoTestSuite))
}
