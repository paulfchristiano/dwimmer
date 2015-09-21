package database_test

import (
	"testing"

	. "github.com/paulfchristiano/dwimmer/storage/database"
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
	collection.Empty("testing")
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
			t.Log(x["value"])
			newVal, err := term.LoadSetting(x["value"])
			if err != nil {
				t.Error(err)
			}
			newID := newVal.ID
			oldID := setting.ID
			if newID != oldID {
				t.Errorf("%v != %v", newVal, setting)
			}
			found++
		}
		if x["key"] == 2 {
			newVal, err := term.LoadC(x["value"])
			if err != nil {
				t.Error(err)
			}
			newID := term.IDC(newVal)
			oldID := term.IDC(cc)
			if newID != oldID {
				t.Errorf("%v != %v", newVal, cc)
			}
			found++
		}
	}
	if found < 2 {
		t.Errorf("found %d < 2 items", found)
	}
	{
		newSetting, err := term.LoadSetting(collection.Get(1))
		if err != nil {
			t.Error(err)
		}
		newID := newSetting.ID
		oldID := setting.ID
		if newID != oldID {
			t.Errorf("%v != %v", newSetting, setting)
		}
	}
	{
		newC, err := term.LoadC(collection.Get(2))
		if err != nil {
			t.Error(err)
		}
		newID := term.IDC(newC)
		oldID := term.IDC(cc)
		if newID != oldID {
			t.Errorf("%v != %v", newC, cc)
		}
	}
}
