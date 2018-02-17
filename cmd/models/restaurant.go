package models

type Restaurant struct {
	ID int
	Name string
	City string
	Rating float32
	Menu []Dish
}
