package mongo

import "gopkg.in/mgo.v2"

var (
	_ DataAccessor = new(DataAccess)
	_ Querier      = new(mgo.Query)
)

type DataAccessor interface {
	Find(interface{}) Querier
	Insert(interface{}) error
}

type Querier interface {
	All(interface{}) error
	One(interface{}) error
}

type DataAccess struct {
	Collection *mgo.Collection
}

func (da *DataAccess) Find(query interface{}) Querier {
	return da.Collection.Find(query)
}

func (da *DataAccess) Insert(object interface{}) error {
	return da.Collection.Insert(object)
}
