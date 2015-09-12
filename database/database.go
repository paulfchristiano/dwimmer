package database

import "gopkg.in/mgo.v2"

var db *mgo.Database

func init() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic("failed to connect to database")
	}
	db = session.DB("dwimmer")
}

func Collection(name string) *mgo.Collection {
	return db.C(name)
}
