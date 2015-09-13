package database

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var db *mgo.Database

type C struct {
	collection *mgo.Collection
}

func init() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic("failed to connect to database")
	}
	db = session.DB("dwimmer")
}

func Collection(name string) C {
	return C{db.C(name)}
}

func (c C) Set(key, value interface{}) {
	c.collection.Upsert(bson.M{"key": key}, bson.M{"$set": bson.M{"value": value}})
}

func (c C) Get(key interface{}) interface{} {
	var holder bson.M
	c.collection.Find(bson.M{"key": key}).One(&holder)
	return holder["value"]
}

func (c C) All() []bson.M {
	result := []bson.M{bson.M{}}
	iter := c.collection.Find(nil).Iter()
	for iter.Next(&result[len(result)-1]) {
		result = append(result, bson.M{})
	}
	return result
}

func (c C) Count() int {
	result, err := c.collection.Count()
	if err != nil {
		panic(err)
	}
	return result
}
