package meta

import (
	"github.com/paulfchristiano/dwimmer/data/core"
	"github.com/paulfchristiano/dwimmer/data/represent"
	"github.com/paulfchristiano/dwimmer/dynamics"
	"github.com/paulfchristiano/dwimmer/term"
)

var (
	SettingForChannel = term.Make("what setting corresponds to the channel []?")
	NotAChannel       = term.Make("the argument is not a channel")
)

func init() {
	s := dynamics.ExpectQuestion(term.InitS(), SettingForChannel, "Q", "c")
	s = dynamics.AddSimple(s, term.ViewS(term.Sr("c")))
	s = s.AppendTemplate(term.Channel{}.Head())
	dynamics.AddNative(s, dynamics.Args1(settingForChannel), "c")
}

func settingForChannel(d dynamics.Dwimmer, context *term.SettingT, channel term.T) term.T {
	c, ok := channel.(term.Channel)
	if !ok {
		return NotAChannel.T()
	}
	return core.Answer.T(represent.SettingT(c.Setting))
}
