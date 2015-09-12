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
	s.AppendTemplate(LastAndInitQ, "l")
	s = dynamics.AddSimple(s, term.ViewS(term.Sr("l")))

	t := s.Copy().AppendTemplate(Snoc, "init", "last")
	t = dynamics.AddSimple(t, term.ReturnS(LastAndInit.S(term.Sr("init"), term.Sr("last"))))

	t = s.Copy().AppendTemplate(Singleton, "x")
	t = dynamics.AddSimple(t, term.ReturnS(LastAndInit.S(term.Sr("x"), Empty.S())))

	t = s.Copy().AppendTemplate(Empty)
	t = dynamics.AddSimple(t, term.ReturnS(IsEmpty.S()))

	t = s.Copy().AppendTemplate(Cons, "first", "rest")
	t = dynamics.AddSimple(t, term.AskS(LastAndInitQ.S(term.Sr("rest"))))
	tt := t.Copy().AppendTemplate(LastAndInit, "last", "init")
	tt = dynamics.AddSimple(tt, term.ReturnS(LastAndInit.S(
		term.Sr("last"),
		Cons.S(term.Sr("first"), term.Sr("init")),
	)))
	tt = t.Copy().AppendTemplate(IsEmpty)
	tt = dynamics.AddSimple(tt, term.ReturnS(LastAndInit.S(term.Sr("first"), Empty.S())))

	t = s.Copy().AppendTemplate(Concat, "a", "b")
	t = dynamics.AddSimple(t, term.AskS(LastAndInitQ.S(term.Sr("b"))))
	tt = t.Copy().AppendTemplate(IsEmpty)
	tt = dynamics.AddSimple(tt, term.AskS(LastAndInitQ.S(term.Sr("a"))))
	dynamics.AddSimple(tt.Copy().AppendTemplate(IsEmpty), term.ReturnS(IsEmpty.S()))
	dynamics.AddSimple(
		tt.Copy().AppendTemplate(LastAndInit, "last", "init"),
		term.ReturnS(LastAndInit.S(term.Sr("last"), term.Sr("init"))),
	)
	tt = t.Copy().AppendTemplate(LastAndInit, "last", "init")
	dynamics.AddSimple(tt, term.ReturnS(LastAndInit.S(
		term.Sr("last"),
		Concat.S(term.Sr("a"), term.Sr("init")),
	)))

}
