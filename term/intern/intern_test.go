package intern

import "testing"

func TestIntern(t *testing.T) {
	strings := []string{
		"hello, world",
		"asdf",
		"testing",
		"asdf",
		"asdf",
		"another test",
	}
	ids := make([]Packed, 50)
	for _, x := range []Packed{new(Id), new(Record)} {
		for i := range ids {
			ids[i] = x
		}
		for i, s := range strings {
			ids[i].StoreStr(s)
			is := ids[i].Str()
			if s != is {
				t.Errorf("%v != %v", s, is)
			}
		}
		for i1, s1 := range strings {
			for i2, s2 := range strings {
				if s1 == s2 && ids[i1] != ids[i2] {
					t.Errorf("%v and %v have ids %v and %v", s1, s2, ids[i1], ids[i2])
				}
			}
		}
		ints := []int{
			1,
			2,
			8,
			123,
			2,
			8,
		}
		for i, n := range ints {
			ids[i].StoreInt(n)
			in := ids[i].Int()
			if n != in {
				t.Errorf("%v != %v", n, in)
			}
		}
		for i1, n1 := range ints {
			for i2, n2 := range ints {
				if n1 == n2 && ids[i1] != ids[i2] {
					t.Errorf("%v and %v have ids %v and %v", n1, n2, ids[i1], ids[i2])
				}
			}
		}
		for i, s := range strings {
			n := ints[i]
			ids[i].StorePair(ids[i].New().StoreInt(n), ids[i].New().StoreStr(s))
			ind, isd := ids[i].Pair()
			in := ind.Int()
			is := isd.Str()
			if n != in {
				t.Errorf("%v != %v", n, in)
			}
			if s != is {
				t.Errorf("%v != %v", s, is)
			}
		}
		for i1, s1 := range strings {
			for i2, s2 := range strings {
				if ints[i1] == ints[i2] && s1 == s2 && ids[i1] != ids[i2] {
					t.Errorf("pairs didn't match up!")
				}
			}
		}
	}
}

func TestLists(t *testing.T) {
	for _, x := range []Packed{new(Id), new(Record)} {
		list := []Packed{
			x.New().StoreInt(1),
			x.New().StoreInt(2),
			x.New().StoreInt(2),
			x.New().StoreInt(3),
			x.New().StoreInt(1),
			x.New().StoreInt(3),
		}
		lid := x.New().StoreList(list)
		list = append(list, x.New().StoreInt(4))
		lid.Append(lid.New().StoreInt(4))
		for j := range []int{0, 1} {
			if j == 1 {
				lid.Empty()
				for _, q := range list {
					lid.Append(q)
				}
			}
			il := lid.List()
			for i, x := range list {
				if il[i].Int() != x.Int() {
					t.Errorf("%v != %v", il[i].Int(), x.Int())
				}
			}
			if len(il) != len(list) {
				t.Errorf("%d != %d", len(il), len(list))
			}
		}
	}
}
