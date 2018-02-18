package models

import "gopkg.in/mgo.v2/bson"

const RestaurantCollectionName = "Restaurants"

type Restaurant struct {
	ID     bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Name   string        `bson:"name,omitempty" json:"name,omitempty" validate:"required"`
	City   string        `bson:"city,omitempty" json:"city,omitempty" query:"city" validate:"required,city"`
	Rating float32       `bson:"rating,omitempty" json:"rating,omitempty" validate:"isdefault=0,min=0,max=10"`
}
