package strings

import (
	"fmt"

	"github.com/paulfchristiano/dwimmer/data/core"
	"github.com/paulfchristiano/dwimmer/data/represent"
	"github.com/paulfchristiano/dwimmer/dynamics"
	"github.com/paulfchristiano/dwimmer/term"
)

func concat(d dynamics.Dwimmer, s *term.SettingT, a, b term.T) term.T {
	return core.Answer.T(term.Str(string(a.(term.Str)) + string(b.(term.Str))))
}

func length(d dynamics.Dwimmer, s *term.SettingT, a term.T) term.T {
	return core.Answer.T(term.Int(len(string(a.(term.Str)))))
}

func bracketed(d dynamics.Dwimmer, s *term.SettingT, a term.T) term.T {
	return core.Answer.T(term.Str(fmt.Sprint("[%s]", string(a.(term.Str)))))
}

var (
	Len       = term.Make("what is the length of []?")
	Concat    = term.Make("what is the result of concatenating [] to []?")
	Bracketed = term.Make("what is the string formed by enclosing [] in brackets?")
	GetChar   = term.Make("what is the character that comes after [] others in the string []?")
)

func init() {
	s := term.InitS()
	s = dynamics.ExpectQuestion(s, Len, "Q", "s")
	s = dynamics.AddSimple(s, term.ViewS(term.Sr("s")))
	s.AppendTemplate(term.Str("").Head())
	dynamics.AddNative(s, dynamics.Args1(length), "s")

	s = term.InitS()
	s = dynamics.ExpectQuestion(s, Bracketed, "Q", "s")
	s = dynamics.AddSimple(s, term.ViewS(term.Sr("s")))
	s.AppendTemplate(term.Str("").Head())
	dynamics.AddNative(s, dynamics.Args1(bracketed), "s")

	s = term.InitS()
	s = dynamics.ExpectQuestion(s, Concat, "Q", "a", "b")
	s = dynamics.AddSimple(s, term.ViewS(term.Sr("a")))
	s.AppendTemplate(term.Str("").Head())
	s = dynamics.AddSimple(s, term.ViewS(term.Sr("b")))
	s.AppendTemplate(term.Str("").Head())
	dynamics.AddNative(s, dynamics.Args2(concat), "a", "b")

	dynamics.AddNativeResponse(GetChar, 2, dynamics.Args2(getChar))
}

func getChar(d dynamics.Dwimmer, s *term.SettingT, quotedN, quotedS term.T) term.T {
	n, err := represent.ToInt(d, quotedN)
	if err != nil {
		return term.Make("asked to index into a string, but received " +
			"[] while converting the index to native format").T(err)
	}
	str, err := represent.ToStr(d, quotedS)
	if err != nil {
		return term.Make("asked to index into a string, but received " +
			"[] while converting the string to native format").T(err)
	}
	return core.Answer.T(represent.Rune(rune(str[n])))
}
