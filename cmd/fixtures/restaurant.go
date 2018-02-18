package fixtures

import (
	"venues/cmd/models"

	"gopkg.in/mgo.v2/bson"
)

func SimpleRestaurantSet() []models.Restaurant {
	return []models.Restaurant{
		{ID: bson.NewObjectId(),
			Name:   "Name1",
			City:   "City1",
			Rating: 4.5,
		},
		{ID: bson.NewObjectId(),
			Name:   "Name2",
			City:   "City2",
			Rating: 5.5,
		},
	}
}
