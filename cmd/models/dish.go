package models

import "gopkg.in/mgo.v2/bson"

const DishCollectionName = "Dishes"

// Price has integer type cause it makes better round control
// Price will represent as price multiplied by 100
type Dish struct {
	ID    bson.ObjectId `bson:"_id,omitempty" json:"ID,omitempty"`
	Name  string        `bson:"name,omitempty" json:"name,omitempty"`
	Price int           `bson:"price,omitempty" json:"price,omitempty"`
}
