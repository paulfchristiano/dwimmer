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

func Explicit(n int) term.T {
	switch {
	case n == 0:
		return Zero.T()
	case n < 0:
		return Negative.T(Explicit(-n))
	case n%2 == 0:
		return Double.T(Explicit(n / 2))
	case n%2 == 1:
		return DoublePlusOne.T(Explicit(n / 2))
	default:
		panic("unreachable")
	}
}

var (
	ExplicitQ = term.Make("what is []? the representation should not involve any wrapped Go objects")
)

func makeExplicit(d dynamics.Dwimmer, s *term.SettingT, n term.T) term.T {
	return core.Answer.T(Explicit(int(n.(term.Int))))
}

func init() {
	s := term.InitS().AppendTemplate(ExplicitQ, "x")
	s = dynamics.AddSimple(s, term.ViewS(term.Sr("x")))

	for _, tm := range []term.TemplateId{DoublePlusOne, Double, Negative} {
		t := s.Copy().AppendTemplate(tm, "y")
		t = dynamics.AddSimple(t, term.AskS(ExplicitQ.S(term.Sr("y"))))
		dynamics.AddSimple(t.AppendTemplate(core.Answer, "A"),
			term.ReturnS(core.Answer.S(tm.S(term.Sr("A")))))
	}
	dynamics.AddSimple(s.Copy().AppendTemplate(Zero), term.ReturnS(core.Answer.S(Zero.S())))

	dynamics.AddNative(s.Copy().AppendTemplate(term.Int(0).Head()), dynamics.Args1(makeExplicit), "x")
}

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

var (
	Plus  = term.Make("what is [] + []?")
	Times = term.Make("what is [] * []?")
	Minus = term.Make("what is [] - []?")
	Less  = term.Make("is [] less than []?")
	Equal = term.Make("is [] equal to []?")
)

func init() {
	QNames := []term.TemplateId{Plus, Times, Minus, Equal, Less}
	QFuncs := [](func(dynamics.Dwimmer, *term.SettingT, ...term.T) term.T){
		dynamics.Args2(addNative),
		dynamics.Args2(multiplyNative),
		dynamics.Args2(subtractNative),
		dynamics.Args2(testEquality),
		dynamics.Args2(testLess),
	}

	for i := range QNames {
		s := term.InitS()
		s.AppendTemplate(QNames[i], "a", "b")
		s = dynamics.AddSimple(s, term.ViewS(term.Sr("a")))
		s.AppendTemplate(term.Int(0).Head())
		s = dynamics.AddSimple(s, term.ViewS(term.Sr("b")))
		s.AppendTemplate(term.Int(0).Head())
		dynamics.AddNative(s, QFuncs[i], "a", "b")
	}
}

func addNative(d dynamics.Dwimmer, s *term.SettingT, n, m term.T) term.T {
	return core.Answer.T(term.Int(int(n.(term.Int)) + int(n.(term.Int))))
}
func subtractNative(d dynamics.Dwimmer, s *term.SettingT, n, m term.T) term.T {
	return core.Answer.T(term.Int(int(n.(term.Int)) - int(n.(term.Int))))
}
func multiplyNative(d dynamics.Dwimmer, s *term.SettingT, n, m term.T) term.T {
	return core.Answer.T(term.Int(int(n.(term.Int)) * int(n.(term.Int))))
}
