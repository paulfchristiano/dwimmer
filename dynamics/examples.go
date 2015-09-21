package dynamics

import "github.com/paulfchristiano/dwimmer/term"

func ExpectQuestion(s *term.SettingS, t term.TemplateID, names ...string) *term.SettingS {
	result := s.Copy()
	result.AppendTemplate(ParentChannel, names[0])
	result.AppendTemplate(t, names[1:]...)
	return result
}

func ExpectAnswer(s *term.SettingS, t term.TemplateID, names ...string) *term.SettingS {
	result := s.Copy()
	result.AppendTemplate(OpenChannel, names[0])
	result.AppendTemplate(t, names[1:]...)
	return result
}

func Parent(s *term.SettingT) term.T {
	return ParentChannel.T(term.MakeChannel(s))
}

var (
	OpenChannel   = term.Make("@[]")
	ParentChannel = term.Make("@[]*")
)
