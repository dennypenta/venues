package repositories

import (
	"venues/cmd/models"
	"venues/pkg/mongo"

	"venues/cmd/storages"

	"gopkg.in/mgo.v2/bson"
)

type RestaurantAccessor interface {
	List() ([]models.Restaurant, error)
}

type RestaurantRepo struct {
	storage mongo.DataAccessor
}

func (repo *RestaurantRepo) List() ([]models.Restaurant, error) {
	var restaurants []models.Restaurant
	err := repo.storage.Find(bson.M{}).All(&restaurants)
	return restaurants, err
}

func NewRestaurantRepo() *RestaurantRepo {
	collection := storages.GetStorage().C(models.RestaurantCollectionName)
	dataAccessor := &mongo.DataAccess{Collection: collection}
	return &RestaurantRepo{storage: dataAccessor}
}