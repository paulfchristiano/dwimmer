package database

import (
	"testing"

	"github.com/paulfchristiano/dwimmer/term"
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
	template := term.Make("test [] term")
	c := template.C(term.ReferenceC{0})
	cc := term.ConstC{template.T(term.Make("stub").T())}
	action := term.ReturnC(c)
	settingS := term.InitS().AppendTemplate(template, "q").AppendAction(action)
	setting := settingS.Setting
	collection.Set(1, term.SaveSetting(setting))
	collection.Set(2, term.SaveC(cc))
	found := 0
	for _, x := range collection.All() {
		if x["key"] == 1 {
			newVal := term.LoadSetting(x["value"])
			newId := newVal.Id
			oldId := setting.Id
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
	{
		newSetting := term.LoadSetting(collection.Get(1))
		newId := newSetting.Id
		oldId := setting.Id
		if newId != oldId {
			t.Errorf("%v != %v", newSetting, setting)
		}
	}
	{
		newC := term.LoadC(collection.Get(2))
		newId := term.IdC(newC)
		oldId := term.IdC(cc)
		if newId != oldId {
			t.Errorf("%v != %v", newC, cc)
		}
	}
}
