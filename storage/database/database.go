package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

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

func (c C) Insert(data bson.M) {
	c.collection.Insert(data)
}

func (c C) LoadFrom(s string) {
	f, err := os.Open(s)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for line, err := r.ReadSlice('\n'); err == nil; line, err = r.ReadSlice('\n') {
		fmt.Println(line)
		var data bson.M
		err = json.Unmarshal(line, &data)
		if err != nil {
			panic(err)
		}
		c.Insert(data)
	}
}

func (c C) DumpTo(s string) {
	f, err := os.Create(s)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	for _, x := range c.All() {
		out, err := json.Marshal(x)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(f, "%s\n", out)
	}
}
