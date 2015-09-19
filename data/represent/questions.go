package represent

import (
	"github.com/paulfchristiano/dwimmer/data/core"
	"github.com/paulfchristiano/dwimmer/dynamics"
	"github.com/paulfchristiano/dwimmer/term"
)

var (
	RepresentSetting = term.Make("what term represents the setting []?")
)

func init() {
	s := term.InitS()
	s = dynamics.ExpectQuestion(s, RepresentSetting, "Q", "s")
	dynamics.AddNative(s, dynamics.Args1(quote), "s")
}

func quote(d dynamics.Dwimmer, s *term.SettingT, t term.T) term.T {
	return core.Answer.T(term.Quote(t))
}
