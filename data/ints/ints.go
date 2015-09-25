package ints

import (
	"github.com/paulfchristiano/dwimmer/data/core"
	"github.com/paulfchristiano/dwimmer/dynamics"
	"github.com/paulfchristiano/dwimmer/term"
)

var (
	DoublePlusOne = term.Make("two times [] plus one")
	Double        = term.Make("two times []")
	Zero          = term.Make("the number zero")
	One           = term.Make("the number one")
	Negative      = term.Make("negative []")
)

/*
var (
	ExplicitQ = term.Make("what is []? the representation should not involve any wrapped Go objects")
)

func makeExplicit(d dynamics.Dwimmer, s *term.SettingT, n term.T) term.T {
	return core.Answer.T(Explicit(int(n.(term.Int))))
}

func init() {
	s := term.InitS().AppendTemplate(ExplicitQ, "x")
	s = dynamics.AddSimple(s, term.ViewS(term.Sr("x")))

	for _, tm := range []term.TemplateID{DoublePlusOne, Double, Negative} {
		t := s.Copy().AppendTemplate(tm, "y")
		t = dynamics.AddSimple(t, term.AskS(ExplicitQ.S(term.Sr("y"))))
		dynamics.AddSimple(t.AppendTemplate(core.Answer, "A"),
			term.ReturnS(core.Answer.S(tm.S(term.Sr("A")))))
	}
	dynamics.AddSimple(s.Copy().AppendTemplate(Zero), term.ReturnS(core.Answer.S(Zero.S())))

	dynamics.AddNative(s.Copy().AppendTemplate(term.Int(0).Head()), dynamics.Args1(makeExplicit), "x")
}
*/

func testEquality(d dynamics.Dwimmer, s *term.SettingT, n, m term.T) term.T {
	if int(n.(term.Int)) == int(m.(term.Int)) {
		return core.Answer.T(core.Yes.T())
	}
	return core.Answer.T(core.No.T())
}

func testLess(d dynamics.Dwimmer, s *term.SettingT, n, m term.T) term.T {
	if int(n.(term.Int)) < int(m.(term.Int)) {
		return core.Answer.T(core.Yes.T())
	}
	return core.Answer.T(core.No.T())
}
func testLessOrEqual(d dynamics.Dwimmer, s *term.SettingT, n, m term.T) term.T {
	if int(n.(term.Int)) <= int(m.(term.Int)) {
		return core.Answer.T(core.Yes.T())
	}
	return core.Answer.T(core.No.T())
}
func testMore(d dynamics.Dwimmer, s *term.SettingT, n, m term.T) term.T {
	if int(n.(term.Int)) > int(m.(term.Int)) {
		return core.Answer.T(core.Yes.T())
	}
	return core.Answer.T(core.No.T())
}
func testMoreOrEqual(d dynamics.Dwimmer, s *term.SettingT, n, m term.T) term.T {
	if int(n.(term.Int)) >= int(m.(term.Int)) {
		return core.Answer.T(core.Yes.T())
	}
	return core.Answer.T(core.No.T())
}

var (
	Plus        = term.Make("what is [] + []?")
	Times       = term.Make("what is [] * []?")
	Minus       = term.Make("what is [] - []?")
	Divide      = term.Make("what is [] / []?")
	Equal       = term.Make("is [] equal to []?")
	Less        = term.Make("is [] less than []?")
	More        = term.Make("is [] greater than []?")
	LessOrEqual = term.Make("is [] at most []?")
	MoreOrEqual = term.Make("is [] at least []?")
)

func init() {
	QNames := []term.TemplateID{Plus, Times, Minus, Equal, Less, More,
		LessOrEqual, MoreOrEqual,
		Divide,
	}
	QFuncs := [](func(dynamics.Dwimmer, *term.SettingT, ...term.T) term.T){
		dynamics.Args2(addNative),
		dynamics.Args2(multiplyNative),
		dynamics.Args2(subtractNative),
		dynamics.Args2(testEquality),
		dynamics.Args2(testLess),
		dynamics.Args2(testMore),
		dynamics.Args2(testLessOrEqual),
		dynamics.Args2(testMoreOrEqual),
		dynamics.Args2(divideNative),
	}

	for i := range QNames {
		s := term.InitS()
		s = dynamics.ExpectQuestion(s, QNames[i], "Q", "a", "b")
		s = dynamics.AddSimple(s, term.ViewS(term.Sr("a")))
		s.AppendTemplate(term.Int(0).Head())
		s = dynamics.AddSimple(s, term.ViewS(term.Sr("b")))
		s.AppendTemplate(term.Int(0).Head())
		dynamics.AddNative(s, QFuncs[i], "a", "b")
	}
}

func addNative(d dynamics.Dwimmer, s *term.SettingT, n, m term.T) term.T {
	return core.Answer.T(term.Int(int(n.(term.Int)) + int(m.(term.Int))))
}
func subtractNative(d dynamics.Dwimmer, s *term.SettingT, n, m term.T) term.T {
	return core.Answer.T(term.Int(int(n.(term.Int)) - int(m.(term.Int))))
}
func multiplyNative(d dynamics.Dwimmer, s *term.SettingT, n, m term.T) term.T {
	return core.Answer.T(term.Int(int(n.(term.Int)) * int(m.(term.Int))))
}
func divideNative(d dynamics.Dwimmer, s *term.SettingT, n, m term.T) term.T {
	a := int(n.(term.Int))
	b := int(m.(term.Int))
	remainder := a % b
	if remainder == 0 {
		return core.Answer.T(term.Int(a / b))
	}
	return QuotientAndRemainder.T(term.Int(a/b), term.Int(remainder))
}

var (
	QuotientAndRemainder = term.Make("the quotient is [] with remainder []")
)
