package models

import "gopkg.in/mgo.v2/bson"

type Menu struct {
	Dishes []Dish `bson: "dishes,omitempty" json:"dishes,omitempty"`
}

// Price has integer type cause it makes better round control
// Price will represent as price multiplied by 100
type Dish struct {
	ID    bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Name  string        `bson:"name,omitempty" json:"name,omitempty"`
	Price int           `bson:"price,omitempty" json:"price,omitempty"`
}

//db.Restaurants.find({'dishes.name': 'name'}, {dishes: 1, _id: 0})
