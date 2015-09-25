package term

import (
	"fmt"
	"testing"
)

/*
func ExampleTemplate() {
	fmt.Println(Make("the pair ([], [])"))
	//Output:
	//the pair ([], [])
}

func ExampleTerm() {
	test := Make("testing []")
	pair := Make("the pair ([], [])")
	stub := Make("test")
	fmt.Println(pair.T(test.T(stub.T()), stub.T()))
	//Output:
	//the pair ([], [])
}

func ExampleSC() {
	test := Make("testing []")
	r := ReferenceS{"testvar"}
	fmt.Println(test.S(r))
	q := ReferenceC{1}
	fmt.Println(test.C(q))
	//Output:
	//testing [#testvar]
	//testing [#1]
}
*/

func TestID(t *testing.T) {
	ts := []T{
		Make("test").T(),
		Int(7),
		Str("hello"),
	}
	for _, t := range ts {
		it := IDT(t).T()
		if it.String() != t.String() {
			fmt.Println("%v != %v", t, it)
		}
		s := ConstC{t}
		is := IDC(s).C()
		if s.String() != is.String() {
			fmt.Println("%v != %v", s, is)
		}
	}
}

func TestActionCID(t *testing.T) {
	cs := []C{Cc(Int(3)), Cr(1), Make("testing []").C(Cr(0))}
	as := []ActionC{
		ReturnC(cs[0]),
		ViewC(cs[1]),
		AskC(cs[2]),
		ClarifyC(cs[2], cs[1]),
		ReplaceC(cs[1], -1),
		CorrectC(3),
	}
	for _, a := range as {
		id := IDActionC(a)
		newa := id.ActionC()
		if a.String() != newa.String() {
			t.Errorf("%v != %v", a, newa)
		}

	}

}

func TestSaving(t *testing.T) {
	temp := Make("a")
	_, ok := LoadTemplate(SaveTemplate(temp.Template()))
	t.Log(SaveTemplate(temp.Template()))
	if !ok {
		t.Errorf("failed to load")
	}
	a := temp.T()
	b := Make("b []").T(a)
	c := Make("c []").T(a)
	b, ok = LoadT(SaveT(b))
	if !ok {
		t.Error("failed to load")
	}
	c, ok = LoadT(SaveT(c))
	if !ok {
		t.Error("failed to load")
	}
	if b.Values()[0] != c.Values()[0] {
		t.Errorf("saving did not collapse instances")
	}
}

func TestCaching(t *testing.T) {
	temp := Make("xyzzy")
	a := Make("like xyzzy []").T(temp.T())
	b := Make("like xyzzy []").T(temp.T())
	if SaveT(a) != SaveT(b) {
		t.Errorf("SaveT(%v) = %v != %v = SaveT(%v)", a, SaveT(a), SaveT(b), b)
	}
}

func BenchmarkSave2(b *testing.B)  { benchmarkSave(b, 2) }
func BenchmarkSave5(b *testing.B)  { benchmarkSave(b, 5) }
func BenchmarkSave10(b *testing.B) { benchmarkSave(b, 10) }

func benchmarkSave(b *testing.B, m int) {
	for iter := 0; iter < b.N; iter++ {
		stub := Make("stub").T()
		pair := Make("the pair ([], [])")
		grid := make([][]T, m)
		for i := range grid {
			grid[i] = make([]T, m)
		}
		for i := range grid {
			for j := range grid[i] {
				if i == 0 || j == 0 {
					grid[i][j] = stub
				} else {
					grid[i][j] = pair.T(grid[i-1][j], grid[i][j-1])
				}
			}
		}
		LoadT(SaveT(grid[m-1][m-1]))
	}
}
