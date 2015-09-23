package intern

import (
	"fmt"
	"testing"

	"github.com/paulfchristiano/dwimmer/storage/database"
)

func TestIntern(t *testing.T) {
	strings := []string{
		"hello, world",
		"asdf",
		"testing",
		"asdf",
		"asdf",
		"another test",
	}
	db := database.Collection("testing")
	packers := []Packer{NewIDer(), NewRecorder(db)}
	for _, packer := range packers {
		ids := make([]Packed, 50)
		for i, s := range strings {
			ids[i] = packer.PackString(s)
			is, ok := packer.UnpackString(ids[i])
			if !ok {
				t.Error(ids[i])
			}
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
			ids[i] = packer.PackInt(n)
			in, ok := packer.UnpackInt(ids[i])
			if !ok {
				t.Error(ids[i])
			}
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
			ids[i] = packer.PackPair(packer.PackInt(n), packer.PackString(s))
			ind, isd, ok := packer.UnpackPair(ids[i])
			if !ok {
				t.Error(ids[i])
			}
			in, ok := packer.UnpackInt(ind)
			if !ok {
				t.Error(ind)
			}
			is, ok := packer.UnpackString(isd)
			if !ok {
				t.Error(isd)
			}
			if n != in {
				t.Errorf("%v != %v", n, in)
			}
			if s != is {
				t.Errorf("%v != %v", s, is)
			}
		}
	}
}

func TestLists(t *testing.T) {
	db := database.Collection("testing")
	for _, packer := range []Packer{NewIDer(), NewRecorder(db)} {
		list := []Packed{
			packer.PackInt(1),
			packer.PackInt(2),
			packer.PackInt(3),
			packer.PackInt(3),
			packer.PackInt(2),
			packer.PackInt(1),
		}
		lid := packer.PackList(list)

		list = append(list, packer.PackInt(4))
		lid = packer.AppendToPacked(lid, packer.PackInt(4))

		il, ok := packer.UnpackList(lid)
		if !ok {
			t.Error(lid)
		}
		fmt.Println(packer.UnpackPair(lid))
		for i, x := range list {
			a, ok := packer.UnpackInt(il[i])
			if !ok {
				t.Error(il[i])
			}
			b, ok := packer.UnpackInt(x)
			if !ok {
				t.Error(x)
			}
			if a != b {
				t.Errorf("%v != %v at %d for %T", a, b, i, packer)
			}
		}
		if len(il) != len(list) {
			t.Errorf("%d != %d", len(il), len(list))
		}
	}
}
