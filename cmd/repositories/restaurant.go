package repositories

import (
	"venues/cmd/models"
	"venues/pkg/mongo"

	"venues/cmd/storages"

	"gopkg.in/mgo.v2/bson"
)

var (
	_ RestaurantAccessor = new(RestaurantRepo)
)

type RestaurantAccessor interface {
	Create(*models.Restaurant) error
	List() ([]models.Restaurant, error)
	Update(*models.Restaurant, *models.Restaurant) error
}

type RestaurantRepo struct {
	storage mongo.DataAccessor
}

func (repo *RestaurantRepo) List() ([]models.Restaurant, error) {
	var restaurants []models.Restaurant
	err := repo.storage.Find(bson.M{}).All(&restaurants)
	return restaurants, err
}

func (repo *RestaurantRepo) Create(object *models.Restaurant) error {
	return repo.storage.Insert(object)
}

func (repo *RestaurantRepo) Update(query *models.Restaurant, object *models.Restaurant) error {
	return repo.storage.Update(query, object)
}

func NewRestaurantRepo() *RestaurantRepo {
	collection := storages.GetStorage().C(models.RestaurantCollectionName)
	dataAccessor := &mongo.DataAccess{Collection: collection}
	return &RestaurantRepo{storage: dataAccessor}
}
