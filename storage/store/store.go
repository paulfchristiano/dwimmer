package store

import (
	"github.com/paulfchristiano/dwimmer/data/core"
	"github.com/paulfchristiano/dwimmer/dynamics"
	"github.com/paulfchristiano/dwimmer/term"
)

var (
	GetState = term.Make("what is the current state of the interpreter?")
	SetState = term.Make("the state of the interpreter should become []")
)

func init() {
	s := dynamics.ExpectQuestion(term.InitS(), GetState, "Q")
	dynamics.AddNative(s, dynamics.Args0(getState))

	s = dynamics.ExpectQuestion(term.InitS(), SetState, "Q", "s")
	dynamics.AddNative(s, dynamics.Args1(setState), "s")
}

func getState(d dynamics.Dwimmer, s *term.SettingT) term.T {
	return core.Answer.T(d.GetStorage())
}

func setState(d dynamics.Dwimmer, s *term.SettingT, t term.T) term.T {
	d.SetStorage(t)
	return core.OK.T()
}
