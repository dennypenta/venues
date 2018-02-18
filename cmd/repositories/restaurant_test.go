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

type MockErrorQuerier struct {
	mock.Mock
}

func (m *MockErrorQuerier) All(model interface{}) error {
	args := m.Called(model)
	return args.Error(0)
}

func (m *MockErrorQuerier) matchedBy(value interface{}) bool {
	return true
}

type MockErrorDataAccess struct {
	mock.Mock
}

func (m *MockErrorDataAccess) Find(query interface{}) mongo.Querier {
	args := m.Called(query)
	return args.Get(0).(mongo.Querier)
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
		storage: &mongo.DataAccess{suite.storage},
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
	mockedStorage := &MockErrorDataAccess{}
	mockedQuerier := &MockErrorQuerier{}
	suite.repo = &RestaurantRepo{storage: mockedStorage}

	mockedStorage.On("Find", bson.M{}).Return(mockedQuerier)
	mockedQuerier.On("All", mock.MatchedBy(mockedQuerier.matchedBy)).Return(errors.New("Mocked error"))

	_, err := suite.repo.List()

	mockedStorage.AssertExpectations(suite.T())
	mockedQuerier.AssertExpectations(suite.T())
	suite.Assertions.Error(err)
}

func TestRestaurantRepoTestSuite(t *testing.T) {
	suite.Run(t, new(RestaurantRepoTestSuite))
}
