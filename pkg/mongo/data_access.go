package mongo

import "gopkg.in/mgo.v2"

var (
	_ DataAccessor = new(DataAccess)
	_ Querier      = new(mgo.Query)
)

type DataAccessor interface {
	Find(interface{}) Querier
}

type Querier interface {
	All(interface{}) error
}

type DataAccess struct {
	Collection *mgo.Collection
}

func (da *DataAccess) Find(query interface{}) Querier {
	return da.Collection.Find(query)
}
