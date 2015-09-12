package term

import (
	"fmt"
	"testing"
)

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
	//the pair ([testing [test]], [test])
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

func TestID(t *testing.T) {
	ts := []T{
		Make("test").T(),
		Int(7),
		Str("hello"),
	}
	for _, t := range ts {
		it := IdT(t).T()
		if it.String() != t.String() {
			fmt.Println("%v != %v", t, it)
		}
		s := ConstC{t}
		is := IdC(s).C()
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
		ClarifyC(cs[2], 2),
		ReplaceC(cs[1], -1),
		CorrectC(3),
	}
	for _, a := range as {
		id := IdActionC(a)
		newa := id.ActionC()
		if a.String() != newa.String() {
			t.Errorf("%v != %v", a, newa)
		}

	}

}
