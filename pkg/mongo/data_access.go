package mongo

import "gopkg.in/mgo.v2"

var (
	_ DataAccessor = new(DataAccess)
	_ Querier      = new(Query)
)

type DataAccessor interface {
	Find(interface{}) Querier
	Insert(interface{}) error
	Update(interface{}, interface{}) error
	Remove(interface{}) error
}

type Querier interface {
	All(interface{}) error
	One(interface{}) error
	Sort(string) Querier
}

type DataAccess struct {
	Collection *mgo.Collection
}

func (da *DataAccess) Find(query interface{}) Querier {
	return &Query{da.Collection.Find(query)}
}

func (da *DataAccess) Insert(object interface{}) error {
	return da.Collection.Insert(object)
}

func (da *DataAccess) Update(query interface{}, object interface{}) error {
	return da.Collection.Update(query, object)
}

func (da *DataAccess) Remove(query interface{}) error {
	return da.Collection.Remove(query)
}

type Query struct {
	query *mgo.Query
}

func(q *Query) All(result interface{}) error {
	return q.query.All(result)
}

func(q *Query) One(result interface{}) error {
	return q.query.One(result)
}

func(q *Query) Sort(ordering string) Querier {
	return &Query{q.query.Sort(ordering)}
}
