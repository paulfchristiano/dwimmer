package lists

import (
	"github.com/paulfchristiano/dwimmer/dynamics"
	"github.com/paulfchristiano/dwimmer/term"
)

var (
	Empty     = term.Make("the empty list")
	Cons      = term.Make("the list with first element [] and following elements []")
	Snoc      = term.Make("the list with first elements [] and last element []")
	Concat    = term.Make("the list formed by concatenating [] and []")
	Singleton = term.Make("the list with the single element []")
)

var (
	LastAndInitQ = term.Make("what is the last element of [], and what are the preceding elements?")
	LastAndInit  = term.Make("the last element is [] and the preceding elements are []")
	IsEmpty      = term.Make("that list is empty")
	RemoveLastN  = term.Make("what list is obtained if the last [] elements are removed from []?")
)

func init() {
	s := term.InitS()
	s = dynamics.ExpectQuestion(s, LastAndInitQ, "Q", "l")
	s = dynamics.AddSimple(s, term.ViewS(term.Sr("l")))

	t := s.Copy().AppendTemplate(Snoc, "init", "last")
	t = dynamics.AddSimple(t, term.ReturnS(LastAndInit.S(term.Sr("init"), term.Sr("last"))))

	t = s.Copy().AppendTemplate(Singleton, "x")
	t = dynamics.AddSimple(t, term.ReturnS(LastAndInit.S(term.Sr("x"), Empty.S())))

	t = s.Copy().AppendTemplate(Empty)
	t = dynamics.AddSimple(t, term.ReturnS(IsEmpty.S()))

	t = s.Copy().AppendTemplate(Cons, "first", "rest")
	t = dynamics.AddSimple(t, term.AskS(LastAndInitQ.S(term.Sr("rest"))))

	tt := dynamics.ExpectAnswer(t.Copy(), LastAndInit, "A", "last", "init")
	tt = dynamics.AddSimple(tt, term.ReturnS(LastAndInit.S(
		term.Sr("last"),
		Cons.S(term.Sr("first"), term.Sr("init")),
	)))
	tt = dynamics.ExpectAnswer(t.Copy(), IsEmpty, "A2")
	tt = dynamics.AddSimple(tt, term.ReturnS(LastAndInit.S(term.Sr("first"), Empty.S())))

	t = s.Copy().AppendTemplate(Concat, "a", "b")
	t = dynamics.AddSimple(t, term.AskS(LastAndInitQ.S(term.Sr("b"))))

	tt = dynamics.ExpectAnswer(t, IsEmpty, "A00")
	tt = dynamics.AddSimple(tt, term.AskS(LastAndInitQ.S(term.Sr("a"))))
	dynamics.AddSimple(dynamics.ExpectAnswer(tt, IsEmpty, "A3"), term.ReturnS(IsEmpty.S()))
	dynamics.AddSimple(
		dynamics.ExpectAnswer(tt, LastAndInit, "A4", "last", "init"),
		term.ReturnS(LastAndInit.S(term.Sr("last"), term.Sr("init"))),
	)

	tt = dynamics.ExpectAnswer(t, LastAndInit, "A7", "last", "init")
	dynamics.AddSimple(tt, term.ReturnS(LastAndInit.S(
		term.Sr("last"),
		Concat.S(term.Sr("a"), term.Sr("init")),
	)))

}
