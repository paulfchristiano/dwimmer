package database

import (
	"testing"

	"github.com/paulfchristiano/dwimmer/term"

	"gopkg.in/mgo.v2/bson"
)

/*
const (
	N = 1e4
)

var result bson.M

func TestMain(m *testing.M) {
	flag.Parse()
	c := Collection("benchmark")
	for j := 0; j < N; j++ {
		c.Insert(bson.M{"test": j})
	}
	os.Exit(m.Run())
}

func BenchmarkSearch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		k := 0
		c := Collection("benchmark")
		iter := c.Find(nil).Iter()
		for iter.Next(&result) {
			k++
		}
		fmt.Println(k)
	}
}
*/

func TestEncoding(t *testing.T) {
	collection := Collection("testing")
	collection.Remove(nil)
	template := term.Make("test [] term")
	c := template.C(term.ReferenceC{0})
	cc := term.ConstC{template.T(term.Make("stub").T())}
	action := term.ReturnC(c)
	settingS := term.InitS().AppendTemplate(template, "q").AppendAction(action)
	setting := settingS.Setting()
	collection.Upsert(
		bson.M{"key": 1},
		bson.M{"$set": bson.M{"value": term.SaveSetting(setting)}},
	)
	collection.Upsert(
		bson.M{"key": 2},
		bson.M{"$set": bson.M{"value": term.SaveC(cc)}},
	)
	var x bson.M
	iter := collection.Find(nil).Iter()
	found := 0
	for iter.Next(&x) {
		if x["key"] == 1 {
			newVal := term.LoadSetting(x["value"])
			newId := term.IdSetting(newVal)
			oldId := term.IdSetting(setting)
			if newId != oldId {
				t.Errorf("%v != %v", newVal, setting)
			}
			found++
		}
		if x["key"] == 2 {
			newVal := term.LoadC(x["value"])
			newId := term.IdC(newVal)
			oldId := term.IdC(cc)
			if newId != oldId {
				t.Errorf("%v != %v", newVal, cc)
			}
			found++
		}
	}
	if found < 2 {
		t.Errorf("found %d < 2 items", found)
	}
}
