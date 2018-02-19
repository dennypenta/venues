package repositories

import (
	"venues/cmd/models"
	"venues/pkg/mongo"

	"venues/cmd/storages"
)

var (
	_ RestaurantAccessor = new(RestaurantRepo)
)

type RestaurantAccessor interface {
	Create(*models.Restaurant) error
	List(*models.Restaurant, string, int) ([]models.Restaurant, error)
	Update(*models.Restaurant, *models.Restaurant) error
	Remove(*models.Restaurant) error
}

type RestaurantRepo struct {
	storage mongo.DataAccessor
}

func (repo *RestaurantRepo) List(filter *models.Restaurant, ordering string, page int) ([]models.Restaurant, error) {
	var restaurants []models.Restaurant

	query := repo.storage.Find(filter)
	if ordering != "" {
		query = query.Sort(ordering)
	}

	if page != 0 {
		query = query.Skip(pageSize * (page - 1)).Limit(pageSize)
	}

	err := query.All(&restaurants)

	return restaurants, err
}

func (repo *RestaurantRepo) Create(object *models.Restaurant) error {
	return repo.storage.Insert(object)
}

func (repo *RestaurantRepo) Update(query *models.Restaurant, object *models.Restaurant) error {
	return repo.storage.Update(query, object)
}

func (repo *RestaurantRepo) Remove(query *models.Restaurant) error {
	return repo.storage.Remove(query)
}

func NewRestaurantRepo() *RestaurantRepo {
	collection := storages.GetStorage().C(models.RestaurantCollectionName)
	dataAccess := &mongo.DataAccess{Collection: collection}
	return &RestaurantRepo{storage: dataAccess}
}
